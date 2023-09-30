package dns

import (
	"fmt"
	"github.com/Dreamacro/clash/component/resolver"
	"github.com/igoogolx/itun2socks/internal/conn"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/pkg/pool"
	D "github.com/miekg/dns"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

type Matcher interface {
	GetDns(question string) (resolver.Resolver, constants.DnsRule)
}

var defaultMatcher Matcher
var mux sync.RWMutex

func UpdateMatcher(m Matcher) {
	mux.Lock()
	defer mux.Unlock()
	defaultMatcher = m
}

func getMatcher() Matcher {
	mux.RLock()
	defer mux.RUnlock()
	return defaultMatcher
}

func HandleDnsConn(conn conn.UdpConn) error {
	var err error
	data := pool.NewBytes(pool.BufSize)
	defer pool.FreeBytes(data)
	_, addr, err := conn.ReadFrom(data)
	dnsMessage := new(D.Msg)
	err = dnsMessage.Unpack(data)
	if err != nil {
		return fmt.Errorf("fail to unpack dns message: err: %v", err)
	}
	res, err := handle(dnsMessage)
	if err != nil {
		return fmt.Errorf("fail to hanlde dns message: err: %v", err)
	}
	resData, err := res.Pack()
	_, err = conn.WriteTo(resData, addr)
	return err
}

type Conn struct {
	remoteAddr chan net.Addr
	written    bool
	read       bool
	data       chan []byte
}

func (d *Conn) WriteTo(data []byte, addr net.Addr) (int, error) {
	if d.written {
		return 0, io.EOF
	}
	dnsMessage := new(D.Msg)
	err := dnsMessage.Unpack(data)
	if err != nil {
		return 0, fmt.Errorf("fail to unpack dns message: err: %v", err)
	}
	res, err := handle(dnsMessage)
	if err != nil {
		return 0, fmt.Errorf("fail to hanlde dns message: err: %v", err)
	}
	resData, err := res.Pack()
	d.data <- resData
	d.remoteAddr <- addr
	d.written = true
	return len(data), err
}

func (d *Conn) ReadFrom(data []byte) (int, net.Addr, error) {
	if d.read {
		return 0, nil, io.EOF
	}
	n := copy(data, <-d.data)
	d.read = true
	return n, <-d.remoteAddr, nil
}

func (d *Conn) Close() error {
	return nil
}

func (d *Conn) SetDeadline(_ time.Time) error {
	return nil
}

func (d *Conn) SetReadDeadline(_ time.Time) error {
	return nil
}

func NewConn() *Conn {
	return &Conn{
		data:       make(chan []byte),
		remoteAddr: make(chan net.Addr),
	}
}

func getDnsQuestion(msg *D.Msg) (string, error) {
	if len(msg.Question) == 0 {
		return "", fmt.Errorf("no dns question")
	}
	name := msg.Question[0].Name
	if strings.HasSuffix(name, ".") {
		name = name[:len(name)-1]
	}
	return name, nil
}

func getResponseIp(msg *D.Msg) []net.IP {
	var ips []net.IP
	for _, a := range msg.Answer {
		if a.Header().Rrtype == D.TypeA {
			ip := net.ParseIP(a.(*D.A).A.String())
			ips = append(ips, ip)
		} else if a.Header().Rrtype == D.TypeAAAA {
			ip := net.ParseIP(a.(*D.AAAA).AAAA.String())
			ips = append(ips, ip)
		} else {
			continue
		}
	}
	return ips
}
