package distribution

import (
	"fmt"
	"github.com/metacubex/mihomo/component/resolver"

	"github.com/igoogolx/itun2socks/internal/cfg/distribution/rule_engine"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/internal/dns"
	"github.com/igoogolx/itun2socks/internal/matcher"
	metaC "github.com/metacubex/mihomo/constant"
)

type SystemProxyConfig struct {
}

func (c SystemProxyConfig) ConnMatcher(metadata *metaC.Metadata, _ rule_engine.Rule) (rule_engine.Rule, error) {

	if metadata.Host != "" {
		var rule, err = matcher.GetRuleEngine().Match(metadata.Host, constants.DomainRuleTypes)
		if err == nil {
			return rule, nil
		}
	}

	if metadata.DstIP.String() != "" {
		rule, err := matcher.GetRuleEngine().Match(metadata.DstIP.String(), constants.IpRuleTypes)
		if err == nil {
			return rule, nil
		}
	}

	return nil, fmt.Errorf("no rule found")

}

func NewSystemProxy() (SystemProxyConfig, error) {
	dns.ResetCache()
	resolver.DefaultResolver = nil
	return SystemProxyConfig{}, nil
}
