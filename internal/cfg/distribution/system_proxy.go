package distribution

import (
	"fmt"
	C "github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/cfg/distribution/rule_engine"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/internal/dns"
	"github.com/igoogolx/itun2socks/internal/matcher"
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
	return SystemProxyConfig{}, nil
}
