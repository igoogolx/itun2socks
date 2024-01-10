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

func (c Config) getIpRuleFromDns(ip string) (constants.Policy, bool) {
	cachedDomain, _, ok := GetCachedDnsItem(ip)
	if ok {
		rule, err := c.RuleEngine.Match(cachedDomain)
		if err == nil {
			return rule.GetPolicy(), true
		}
	}
	return constants.PolicyDirect, false
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

	dnsResult, dnsRuleOk := c.getIpRuleFromDns(ip)

	if dnsRuleOk {
		result = dnsResult
	} else {
		rule, err := c.RuleEngine.Match(ip)
		if err == nil {
			result = rule.GetPolicy()
		}
	}

	return result, nil
}

func (c Config) GetDnsType(domain string) (constants.DnsType, error) {
	var rule, err = c.RuleEngine.Match(domain)
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
