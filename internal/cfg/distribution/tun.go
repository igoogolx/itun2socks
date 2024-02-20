package distribution

import (
	"fmt"
	C "github.com/igoogolx/clash/constant"
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

func (c Config) connMatcher(metadata *C.Metadata, prevRule constants.Policy) (constants.Policy, error) {
	ip := metadata.DstIP.String()

	if dnsResult, dnsRuleOk := c.getIpRuleFromDns(ip); dnsRuleOk {
		return dnsResult, nil
	}
	if rule, err := c.RuleEngine.Match(ip); err == nil {
		return rule.GetPolicy(), nil
	}
	return constants.PolicyProxy, nil
}

func (c Config) ConnMatcher(metadata *C.Metadata, prevRule constants.Policy) (constants.Policy, error) {
	result, err := c.connMatcher(metadata, prevRule)
	ip := metadata.DstIP.String()
	domain := "unknown"
	dnsRule := "unknown"
	cacheDomain, cachedRule, ok := GetCachedDnsItem(ip)
	if ok {
		domain = cacheDomain
		dnsRule = string(cachedRule)
	}
	log.Infoln(log.FormatLog(log.RulePrefix, "ip:%v, rule:%v; domain:%v, rule:%v"), ip, result, domain, dnsRule)
	return result, err
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
