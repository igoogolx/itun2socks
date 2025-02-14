package dns

import (
	"context"
	"fmt"
	cResolver "github.com/Dreamacro/clash/component/resolver"
	"github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/cfg/distribution/ruleEngine"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/internal/matcher"
	"github.com/igoogolx/itun2socks/pkg/log"
	D "github.com/miekg/dns"
	"net"
	"strings"
	"sync"
	"time"
)

var dnsMap = map[constants.Policy]cResolver.Resolver{}
var mux sync.RWMutex

func UpdateDnsMap(local, remote cResolver.Resolver) {
	mux.Lock()
	defer mux.Unlock()
	dnsMap[constants.PolicyDirect] = local
	dnsMap[constants.PolicyProxy] = remote
}

type Conn interface {
	ReadFrom([]byte) (int, net.Addr, error)
	WriteTo([]byte, net.Addr) (int, error)
}

func handleMsgWithEmptyAnswer(r *D.Msg) *D.Msg {
	msg := &D.Msg{}
	msg.Answer = []D.RR{}

	msg.SetRcode(r, D.RcodeSuccess)
	msg.Authoritative = true
	msg.RecursionAvailable = true

	return msg
}

func getDnsQuestion(msg *D.Msg) (string, uint16, error) {
	if len(msg.Question) == 0 {
		return "", D.TypeNone, fmt.Errorf("no dns question")
	}
	name := msg.Question[0].Name
	if strings.HasSuffix(name, ".") {
		name = name[:len(name)-1]
	}
	qType := msg.Question[0].Qtype
	return name, qType, nil
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

func convertRulePolicyToResolver(rule ruleEngine.Rule) (ruleEngine.Rule, error) {
	if rule.GetPolicy() == constants.PolicyReject {
		return nil, fmt.Errorf("reject dns")
	}
	return rule, nil
}

func getDnsResolver(domain string, metadata *constant.Metadata, curRuleEngine *ruleEngine.Engine) (ruleEngine.Rule, error) {
	processPath := metadata.ProcessPath
	var rule ruleEngine.Rule
	var err error
	if len(processPath) != 0 {
		rule, err = curRuleEngine.Match(processPath, constants.ProcessRuleTypes)
		if err == nil {
			return convertRulePolicyToResolver(rule)
		}
	}

	rule, err = curRuleEngine.Match(domain, constants.DomainRuleTypes)
	if err == nil {
		return convertRulePolicyToResolver(rule)
	}

	return ruleEngine.BuiltInProxyRule, nil
}

func Handle(dnsMessage *D.Msg, metadata *constant.Metadata) (*D.Msg, error) {
	mux.RLock()
	defer mux.RUnlock()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var err error
	start := time.Now()
	question, qType, err := getDnsQuestion(dnsMessage)

	if err != nil {
		return nil, fmt.Errorf("invalid dns question, err: %v", err)
	}

	switch qType {
	case D.TypeAAAA, D.TypeSVCB, D.TypeHTTPS:
		return handleMsgWithEmptyAnswer(dnsMessage), nil
	}

	var curRuleEngine = matcher.GetRuleEngine()
	dnsRule, err := getDnsResolver(question, metadata, curRuleEngine)
	if err != nil {
		return nil, fmt.Errorf("fail to get dns resolver, err: %v, question: %v", err, question)
	}

	defer func() {
		if err != nil {
			countFailQuery(dnsRule.GetPolicy())
		} else {
			countSuccessQuery(dnsRule.GetPolicy())
		}
	}()

	res, err := dnsMap[dnsRule.GetPolicy()].ExchangeContext(ctx, dnsMessage)
	if err != nil {
		return nil, fmt.Errorf("fail to exchange dns message, err: %v, question: %v", err, question)
	}
	resIps := getResponseIp(res)
	for _, resIp := range resIps {
		if resIp != nil {
			log.Debugln(log.FormatLog(log.DnsPrefix, "add cache, resIp:%v, question: %v, rule: %v"), resIp, question, dnsRule.GetPolicy())
			addCachedDnsItem(resIp.String(), question)
			curRuleEngine.AddCache(resIp.String(), dnsRule)
		}
	}
	elapsed := time.Since(start).Milliseconds()
	log.Infoln(log.FormatLog(log.DnsPrefix, "target: %v, type: %v, time: %v ms, result: %v"), question, D.TypeToString[qType], elapsed, resIps)
	return res, err
}
