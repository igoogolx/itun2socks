package conn

import (
	"fmt"
	"github.com/metacubex/mihomo/adapter"
	"github.com/metacubex/mihomo/adapter/outbound"
	"strings"
	"sync"

	"github.com/igoogolx/itun2socks/internal/cfg/distribution/rule_engine"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/internal/dns"
	"github.com/igoogolx/itun2socks/pkg/log"
	metaC "github.com/metacubex/mihomo/constant"
)

var defaultIsFakeIpEnabled bool

func UpdateIsFakeIpEnabled(value bool) {
	mux.Lock()
	defer mux.Unlock()
	defaultIsFakeIpEnabled = value
}

var (
	proxies map[constants.Policy]metaC.Proxy
	mux     sync.RWMutex
)

type Matcher func(metadata *metaC.Metadata, rule rule_engine.Rule) (rule_engine.Rule, error)

func RejectQuicMather(metadata *metaC.Metadata, prevRule rule_engine.Rule) (rule_engine.Rule, error) {
	if prevRule.GetPolicy() == constants.PolicyProxy && strings.Contains(metadata.NetWork.String(), "udp") && metadata.DstPort == 443 {
		log.Debugln("reject quic conn:%v", metadata.RemoteAddress())
		return rule_engine.BuiltInRejectRule, nil
	}
	return nil, fmt.Errorf("not quic")
}

func UpdateProxy(remoteProxy metaC.Proxy) {
	mux.Lock()
	defer mux.Unlock()
	proxies = make(map[constants.Policy]metaC.Proxy)
	proxies[constants.PolicyProxy] = remoteProxy
	proxies[constants.PolicyDirect] = adapter.NewProxy(outbound.NewDirect())
	proxies[constants.PolicyReject] = adapter.NewProxy(outbound.NewReject())
}

func GetProxy(rule constants.Policy) (metaC.Proxy, error) {
	mux.RLock()
	defer mux.RUnlock()
	connDialer := proxies[rule]
	if connDialer == nil {
		return nil, fmt.Errorf("empty dialer")
	}
	return connDialer, nil
}

func handleMetadata(metadata *metaC.Metadata) rule_engine.Rule {

	rule := resolveMetadata(metadata)

	if rule.Type() == constants.RuleDnsMap {
		dnsMapRule, ok := rule.(*rule_engine.DnsMap)
		if ok {
			ip, err := dnsMapRule.GetIp()
			if err == nil {
				log.Infoln(log.FormatLog(log.DnsPrefix, "query dns from DNS-MAP rule, question: %v, result: %v"), metadata.Host, ip.String())
				metadata.DstIP = ip
				metadata.Host = ""
			}
		}
	}

	if rule.GetPolicy() == constants.PolicyProxy && defaultIsFakeIpEnabled {
		hostByFakeIp, ok := dns.FakeIpPool.LookBack(metadata.DstIP)
		if ok {
			metadata.Host = hostByFakeIp
		}
	}
	return rule
}
