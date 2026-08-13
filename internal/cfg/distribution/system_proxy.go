package distribution

import (
	"fmt"

	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/internal/dns"
	"github.com/igoogolx/itun2socks/internal/matcher"
	cResolver "github.com/igoogolx/itun2socks/pkg/clash/component/resolver"
	C "github.com/igoogolx/itun2socks/pkg/clash/constant"
	"github.com/igoogolx/itun2socks/pkg/rule_engine"
)

type SystemProxyConfig struct {
}

func (c SystemProxyConfig) ConnMatcher(metadata *C.Metadata, _ rule_engine.Rule) (rule_engine.Rule, error) {

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
	cResolver.DefaultResolver = nil
	return SystemProxyConfig{}, nil
}
