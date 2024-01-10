package distribution

import (
	"fmt"
	C "github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/cfg/distribution/ruleEngine"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/pkg/log"
)

type SystemProxyConfig struct {
	RuleEngine *ruleEngine.Engine
}

func (c SystemProxyConfig) ConnMatcher(metadata *C.Metadata, prevRule constants.Policy) (constants.Policy, error) {
	if metadata.Host != "" {
		dnsRule, err := c.GetDnsType(metadata.Host)
		if err != nil {
			return constants.PolicyProxy, fmt.Errorf("reject dns")
		}
		if dnsRule == constants.LocalDns {
			return constants.PolicyDirect, nil
		} else {
			return constants.PolicyProxy, nil
		}
	}

	ip := metadata.DstIP.String()
	result := constants.PolicyProxy

	defer func() {
		domain := metadata.String()
		log.Infoln(log.FormatLog(log.RulePrefix, "host: %v, rule: %v"), domain, result)
	}()

	rule, err := c.RuleEngine.Match(ip)
	if err == nil {
		return rule.GetPolicy(), nil
	}

	return constants.PolicyProxy, nil

}

func (c SystemProxyConfig) GetDnsType(domain string) (constants.DnsType, error) {
	var rule, err = c.RuleEngine.Match(domain)
	if err != nil {
		return constants.RemoteDns, err
	}
	if rule.GetPolicy() == constants.PolicyDirect {
		return constants.LocalDns, nil
	} else if rule.GetPolicy() == constants.PolicyReject {
		return constants.RemoteDns, fmt.Errorf("reject dns")
	}
	return constants.RemoteDns, fmt.Errorf("dns type not found")
}

func NewSystemProxy(
	ruleId string,
	rules []string,
) (SystemProxyConfig, error) {
	ResetCache()
	rEngine, err := ruleEngine.New(ruleId, rules)
	if err != nil {
		return SystemProxyConfig{}, err
	}
	return SystemProxyConfig{
		RuleEngine: rEngine,
	}, nil
}
