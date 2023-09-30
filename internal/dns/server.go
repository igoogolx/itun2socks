package dns

import (
	"context"
	"fmt"
	"github.com/igoogolx/itun2socks/internal/cfg/distribution"
	"github.com/igoogolx/itun2socks/pkg/log"
	"github.com/miekg/dns"
	"time"
)

var server *dns.Server

func Start() {
	server = &dns.Server{Addr: ":53", Net: "udp"}
	dns.HandleFunc(".", handleDNSRequest)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalln("Failed to start DNS server: %v\n", err)
		}
	}()
}

func Shutdown() {
	err := server.Shutdown()
	if err != nil {
		log.Errorln("Failed to shutdown dns server: %v", err)
	}
}

func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	m, err := handle(r)
	if err != nil {
		log.Errorln("Failed to handle DNS: %v\n", err)
		return
	}
	err = w.WriteMsg(m)
	if err != nil {
		log.Errorln("Failed to write dns message: %v\n", err)
		return
	}
	log.Infoln("write dns message back")
}

func handle(dnsMessage *dns.Msg) (*dns.Msg, error) {
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
	dnsClient, dnsRule := getMatcher().GetDns(question)
	res, err := dnsClient.ExchangeContext(ctx, dnsMessage)
	if err != nil {
		return nil, fmt.Errorf("fail to exchange dns message, err: %v, question: %v", err, question)
	}
	resIps := getResponseIp(res)
	for _, resIp := range resIps {
		if resIp != nil {
			log.Debugln(log.FormatLog(log.DnsPrefix, "add cache, resIp:%v, question: %v, rule: %v"), resIp, question, dnsRule)
			distribution.AddCachedDnsItem(resIp.String(), question, dnsRule)
		}
	}
	log.Infoln(log.FormatLog(log.DnsPrefix, "target: %v, result: %v"), question, resIps)
	return res, err
}
