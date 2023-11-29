package distribution

import (
	C "github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/cfg/distribution/ruleEngine"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/pkg/log"
)

type SystemProxyConfig struct {
	RuleEngine *ruleEngine.Engine
}

func (c SystemProxyConfig) GetConnRule(metadata C.Metadata) constants.RuleType {
	if metadata.Host != "" {
		dnsRule := c.GetDnsType(metadata.Host)
		if dnsRule == constants.LocalDns {
			return constants.RuleBypass
		} else {
			return constants.RuleProxy
		}
	}

	ip := metadata.DstIP.String()
	result := constants.RuleProxy

	defer func() {
		domain := metadata.String()
		log.Infoln(log.FormatLog(log.RulePrefix, "host: %v, rule: %v"), domain, result)
	}()

	rule, err := c.RuleEngine.Match(ip)
	if err == nil {
		return rule.GetPolicy()
	}

	return constants.RuleProxy

}

func (c SystemProxyConfig) GetDnsType(domain string) constants.DnsType {
	var rule, err = c.RuleEngine.Match(domain)
	if err == nil {
		if rule.GetPolicy() == constants.RuleBypass {
			return constants.LocalDns
		} else if rule.GetPolicy() == constants.RuleProxy {
			return constants.RemoteDns
		}
	}
	return constants.RemoteDns
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
