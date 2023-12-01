package distribution

import (
	"fmt"
	C "github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/cfg/distribution/ruleEngine"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/pkg/log"
	"strings"
)

type Config struct {
	Dns        DnsDistribution
	RuleEngine *ruleEngine.Engine
}

func (c Config) getIpRuleFromDns(ip string) (constants.RuleType, error) {
	result := constants.RuleBypass
	_, cachedRule, ok := GetCachedDnsItem(ip)
	if ok {
		if cachedRule == constants.LocalDns {
			result = constants.RuleBypass
		} else {
			result = constants.RuleProxy
		}
		return result, nil
	}
	return constants.RuleBypass, fmt.Errorf("not found")
}

func (c Config) ConnMatcher(metadata *C.Metadata) (constants.RuleType, error) {
	ip := metadata.DstIP.String()

	result := constants.RuleProxy

	defer func(latestIp string) {
		domain := "unknown"
		dnsRule := "unknown"
		cacheDomain, cachedRule, ok := GetCachedDnsItem(latestIp)
		if ok {
			domain = cacheDomain
			dnsRule = string(cachedRule)
		}
		log.Infoln(log.FormatLog(log.RulePrefix, "ip:%v, rule:%v; domain:%v, rule:%v"), latestIp, result, domain, dnsRule)
	}(ip)

	var err error
	//dns result
	result, err = c.getIpRuleFromDns(ip)
	if err == nil {
		return result, nil
	}

	rule, err := c.RuleEngine.Match(ip)
	if err == nil {
		return rule.GetPolicy(), nil
	}
	return constants.RuleProxy, nil

}

func (c Config) GetDnsTypeFromRuleEngine(domain string) (constants.DnsType, error) {
	var rule, err = c.RuleEngine.Match(domain)
	if err != nil {
		return constants.LocalDns, err
	}
	if rule.GetPolicy() == constants.RuleBypass {
		return constants.LocalDns, nil
	} else if rule.GetPolicy() == constants.RuleProxy {
		return constants.RemoteDns, nil
	}
	return constants.LocalDns, fmt.Errorf("dns rule not found")
}

func (c Config) GetDnsType(domain string) constants.DnsType {
	result := constants.RemoteDns
	if strings.Contains(c.Dns.Remote.Address, domain) {
		result = constants.BoostDns
	} else {
		var rule, err = c.GetDnsTypeFromRuleEngine(domain)
		if err == nil {
			result = rule
		}
	}
	return result
}

func NewTun(
	boostDns string,
	remoteDns string,
	localDns string,
	ruleId string,
	rules []string,
	defaultInterfaceName string,
) (Config, error) {
	ResetCache()
	rEngine, err := ruleEngine.New(ruleId, rules)
	if err != nil {
		return Config{}, err
	}
	dns, err := NewDnsDistribution(boostDns, remoteDns, localDns, defaultInterfaceName)
	if err != nil {
		return Config{}, err
	}
	return Config{
		Dns: dns, RuleEngine: rEngine,
	}, nil
}
