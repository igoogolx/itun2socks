package distribution

import (
	"fmt"
	C "github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/cfg/distribution/ruleEngine"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/internal/dns"
	"github.com/igoogolx/itun2socks/internal/matcher"
)

type SystemProxyConfig struct {
}

func (c SystemProxyConfig) ConnMatcher(metadata *C.Metadata, _ ruleEngine.Rule) (ruleEngine.Rule, error) {

	if metadata.Host != "" {
		var rule, err = matcher.GetRule().Match(metadata.Host, constants.DomainRuleTypes)
		if err == nil {
			return rule, nil
		}
	}

	if metadata.DstIP.String() != "" {
		rule, err := matcher.GetRule().Match(metadata.DstIP.String(), constants.IpRuleTypes)
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
