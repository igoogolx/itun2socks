package conn

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Dreamacro/clash/adapter"
	"github.com/Dreamacro/clash/adapter/outbound"
	C "github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/cfg/distribution/rule_engine"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/internal/dns"
	"github.com/igoogolx/itun2socks/pkg/log"
)

var defaultIsFakeIpEnabled bool

func UpdateIsFakeIpEnabled(value bool) {
	mux.Lock()
	defer mux.Unlock()
	defaultIsFakeIpEnabled = value
}

var (
	proxies map[constants.Policy]C.Proxy
	mux     sync.RWMutex
)

type Matcher func(metadata *C.Metadata, rule rule_engine.Rule) (rule_engine.Rule, error)

func RejectQuicMather(metadata *C.Metadata, prevRule rule_engine.Rule) (rule_engine.Rule, error) {
	if prevRule.GetPolicy() == constants.PolicyProxy && strings.Contains(metadata.NetWork.String(), "udp") && metadata.DstPort.String() == "443" {
		log.Debugln("reject quic conn:%v", metadata.RemoteAddress())
		return rule_engine.BuiltInRejectRule, nil
	}
	return nil, fmt.Errorf("not quic")
}

func UpdateProxy(remoteProxy C.Proxy) {
	mux.Lock()
	defer mux.Unlock()
	proxies = make(map[constants.Policy]C.Proxy)
	proxies[constants.PolicyProxy] = remoteProxy
	proxies[constants.PolicyDirect] = adapter.NewProxy(outbound.NewDirect())
	proxies[constants.PolicyReject] = adapter.NewProxy(outbound.NewReject())
}

func GetProxy(rule constants.Policy) (C.Proxy, error) {
	mux.RLock()
	defer mux.RUnlock()
	connDialer := proxies[rule]
	if connDialer == nil {
		return nil, fmt.Errorf("empty dialer")
	}
	return connDialer, nil
}

func handleMetadata(metadata *C.Metadata) rule_engine.Rule {

	rule := resolveMetadata(metadata)

	if rule.Type() == constants.RuleDnsMap {
		dnsMapRule, ok := rule.(*rule_engine.DnsMap)
		if ok {
			ip, err := dnsMapRule.GetIp()
			if err == nil {
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
