package distribution

import (
	"fmt"
	C "github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/cfg/distribution/ruleEngine"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/internal/matcher"
	"github.com/igoogolx/itun2socks/pkg/log"
)

type SystemProxyConfig struct {
}

func (c SystemProxyConfig) connMatcher(metadata *C.Metadata, _ ruleEngine.Rule) (ruleEngine.Rule, error) {

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

func NewSystemProxy() (SystemProxyConfig, error) {
	ResetCache()
	return SystemProxyConfig{}, nil
}
