package dns

import (
	"context"
	"fmt"
	cResolver "github.com/Dreamacro/clash/component/resolver"
	"github.com/igoogolx/itun2socks/internal/cfg/distribution"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/internal/matcher"
	"github.com/igoogolx/itun2socks/pkg/log"
	"github.com/igoogolx/itun2socks/pkg/pool"
	D "github.com/miekg/dns"
	"net"
	"strings"
	"time"
)

var dnsMap = map[constants.DnsType]cResolver.Resolver{}

func UpdateDnsMap(local, remote, boost cResolver.Resolver) {
	dnsMap[constants.LocalDns] = local
	dnsMap[constants.RemoteDns] = remote
	dnsMap[constants.BoostDns] = boost
}

type Conn interface {
	ReadFrom([]byte) (int, net.Addr, error)
	WriteTo([]byte, net.Addr) (int, error)
}

func HandleDnsConn(conn Conn) error {
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

func handle(dnsMessage *D.Msg) (*D.Msg, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	start := time.Now()
	question, err := getDnsQuestion(dnsMessage)
	defer func() {
		elapsed := time.Since(start).Milliseconds()
		log.Debugln(log.FormatLog(log.DnsPrefix, "it took %v ms to handle dns, question: %v"), elapsed, question)
	}()
	if err != nil {
		return nil, fmt.Errorf("invalid dns question, err: %v", err)
	}
	dnsType := matcher.GetDnsMatcher().GetDnsType(question)
	res, err := dnsMap[dnsType].ExchangeContext(ctx, dnsMessage)
	if err != nil {
		return nil, fmt.Errorf("fail to exchange dns message, err: %v, question: %v", err, question)
	}
	resIps := getResponseIp(res)
	for _, resIp := range resIps {
		if resIp != nil {
			log.Debugln(log.FormatLog(log.DnsPrefix, "add cache, resIp:%v, question: %v, rule: %v"), resIp, question, dnsType)
			distribution.AddCachedDnsItem(resIp.String(), question, dnsType)
		}
	}
	log.Infoln(log.FormatLog(log.DnsPrefix, "target: %v, result: %v"), question, resIps)
	return res, err
}
