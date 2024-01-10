package distribution

import (
	"fmt"
	C "github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/cfg/distribution/ruleEngine"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/pkg/log"
)

type Config struct {
	Dns        DnsDistribution
	RuleEngine *ruleEngine.Engine
}

func (c Config) getIpRuleFromDns(ip string) (constants.Policy, error) {
	cachedDomain, _, ok := GetCachedDnsItem(ip)
	if ok {
		rule, err := c.RuleEngine.Match(cachedDomain)
		if err != nil {
			return constants.PolicyProxy, err
		}
		return rule.GetPolicy(), nil
	}
	return constants.PolicyDirect, fmt.Errorf("not found")
}

func (c Config) ConnMatcher(metadata *C.Metadata, prevRule constants.Policy) (constants.Policy, error) {
	ip := metadata.DstIP.String()

	result := constants.PolicyProxy

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
	return constants.PolicyProxy, nil

}

func (c Config) GetDnsTypeFromRuleEngine(domain string) (constants.DnsType, error) {
	var rule, err = c.RuleEngine.Match(domain)
	if err != nil {
		return constants.LocalDns, err
	}
	if rule.GetPolicy() == constants.PolicyDirect {
		return constants.LocalDns, nil
	} else if rule.GetPolicy() == constants.PolicyProxy {
		return constants.RemoteDns, nil
	} else if rule.GetPolicy() == constants.PolicyReject {
		return constants.LocalDns, fmt.Errorf("reject dns")
	}
	return constants.LocalDns, fmt.Errorf("dns rule not found")
}

func (c Config) GetDnsType(domain string) (constants.DnsType, error) {
	return c.GetDnsTypeFromRuleEngine(domain)
}

func NewTun(
	boostDns []string,
	remoteDns []string,
	localDns []string,
	ruleId string,
	rules []string,
	defaultInterfaceName string,
	disableCache bool,
) (Config, error) {
	if len(boostDns) == 0 || len(remoteDns) == 0 || len(localDns) == 0 {
		return Config{}, fmt.Errorf("dns can't be empty")
	}
	ResetCache()
	rEngine, err := ruleEngine.New(ruleId, rules)
	if err != nil {
		return Config{}, err
	}
	dns, err := NewDnsDistribution(boostDns, remoteDns, localDns, defaultInterfaceName, disableCache)
	if err != nil {
		return Config{}, err
	}
	return Config{
		Dns: dns, RuleEngine: rEngine,
	}, nil
}
