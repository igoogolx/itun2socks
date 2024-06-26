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

func (c SystemProxyConfig) connMatcher(metadata *C.Metadata, _ ruleEngine.Rule) (ruleEngine.Rule, error) {

	if metadata.Host != "" {
		var rule, err = c.RuleEngine.Match(metadata.Host, constants.DomainRuleTypes)
		if err == nil {
			return rule, nil
		}
	}

	if metadata.DstIP.String() != "" {
		rule, err := c.RuleEngine.Match(metadata.DstIP.String(), constants.IpRuleTypes)
		if err == nil {
			return rule, nil
		}
	}

	return nil, fmt.Errorf("no rule found")

}

func (c SystemProxyConfig) ConnMatcher(metadata *C.Metadata, prevRule ruleEngine.Rule) (ruleEngine.Rule, error) {
	result, err := c.connMatcher(metadata, prevRule)
	if err != nil {
		return result, err
	}
	defer func() {
		target := metadata.String()
		log.Infoln(log.FormatLog(log.RulePrefix, "host: %v, rule: %v"), target, result.GetPolicy())
	}()
	return result, nil
}

func (c SystemProxyConfig) GetDnsType(domain string, _ *C.Metadata) (constants.DnsType, error) {
	var rule, err = c.RuleEngine.Match(domain, constants.DomainRuleTypes)
	if err == nil {
		if rule.GetPolicy() == constants.PolicyDirect {
			return constants.LocalDns, nil
		} else if rule.GetPolicy() == constants.PolicyProxy {
			return constants.RemoteDns, nil
		} else if rule.GetPolicy() == constants.PolicyReject {
			return constants.LocalDns, fmt.Errorf("reject dns")
		}
	}
	return constants.RemoteDns, nil
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
