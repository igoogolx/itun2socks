package dns

import (
	"context"
	"fmt"
	"github.com/Dreamacro/clash/log"
	"github.com/igoogolx/itun2socks/cfg"
	"github.com/igoogolx/itun2socks/components/resolver"
	D "github.com/miekg/dns"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

type Matcher interface {
	GetDns(question string, isPrimary bool) resolver.Client
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

type Conn struct {
	remoteAddr chan net.Addr
	written    *bool
	read       *bool
	data       chan []byte
	proxyAddr  string
}

func (d Conn) WriteTo(data []byte, addr net.Addr) (int, error) {
	if *d.written {
		return 0, io.EOF
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	start := time.Now()
	dnsMessage := new(D.Msg)
	err := dnsMessage.Unpack(data)
	if err != nil {
		return 0, fmt.Errorf("fail to unpack dns message: err: %v", err)
	}
	question, err := getDnsQuestion(dnsMessage)
	defer func() {
		elapsed := time.Since(start).Milliseconds()
		log.Debugln("[DNS], it took %v ms to handle dns, question: %v", elapsed, question)
	}()
	if err != nil {
		return 0, fmt.Errorf("invalid dns question, err: %v", err)
	}
	dnsClient := getMatcher().GetDns(question, strings.Contains(question, d.proxyAddr))
	res, err := dnsClient.ExchangeContext(ctx, dnsMessage)
	if err != nil {
		return 0, fmt.Errorf("fail to exchage dns message, err: %v, quesion: %v, proxy addr: %v, server: %v", err, question, d.proxyAddr, dnsClient.Nameservers())
	}
	resData, err := res.Pack()
	if err != nil {
		return 0, fmt.Errorf("fail to pack dns responsed, err: %v", err)
	}
	resIps := getResponseIp(res)
	for _, resIp := range resIps {
		if resIp != nil {
			log.Infoln("[DNS], target:%v, server:%v, result: %v", question, dnsClient.Nameservers(), resIp)
			cfg.DnsTable.Add(resIp.String(), question)
		}
	}
	d.data <- resData
	d.remoteAddr <- addr
	*d.written = true
	return len(data), err
}

func (d Conn) ReadFrom(data []byte) (int, net.Addr, error) {
	if *d.read {
		return 0, nil, io.EOF
	}
	n := copy(data, <-d.data)
	*d.read = true
	return n, <-d.remoteAddr, nil
}

func (d Conn) Close() error {

	return nil
}

func (d Conn) SetDeadline(t time.Time) error {

	return nil
}

func (d Conn) SetReadDeadline(t time.Time) error {

	return nil
}

func NewConn(proxyAddr string) Conn {
	done := false
	return Conn{
		data:       make(chan []byte),
		remoteAddr: make(chan net.Addr),
		written:    &done,
		read:       &done,
		proxyAddr:  proxyAddr,
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
