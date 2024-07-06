package distribution

import (
	"fmt"
	C "github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/cfg/distribution/ruleEngine"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/internal/matcher"
	"github.com/igoogolx/itun2socks/pkg/log"
)

type Config struct {
	Dns DnsDistribution
}

func (c Config) getIpRuleFromDns(ip string) (ruleEngine.Rule, bool) {
	_, cachedDomainRule, ok := GetCachedDnsItem(ip)
	if ok {
		return cachedDomainRule, true
	}
	return nil, false
}

func (c Config) connMatcher(metadata *C.Metadata, _ ruleEngine.Rule) (ruleEngine.Rule, error) {
	processPath := metadata.ProcessPath
	if len(processPath) != 0 {
		if rule, err := matcher.GetRule().Match(processPath, constants.ProcessRuleTypes); err == nil {
			return rule, nil
		}
	}

	ip := metadata.DstIP.String()
	if dnsResult, dnsRuleOk := c.getIpRuleFromDns(ip); dnsRuleOk {
		return dnsResult, nil
	}
	if rule, err := matcher.GetRule().Match(ip, constants.IpRuleTypes); err == nil {
		return rule, nil
	}
	return nil, fmt.Errorf("no rule found")
}

func (c Config) ConnMatcher(metadata *C.Metadata, prevRule ruleEngine.Rule) (ruleEngine.Rule, error) {
	result, err := c.connMatcher(metadata, prevRule)
	if err != nil {
		return result, err
	}
	ip := metadata.DstIP.String()
	domain := "unknown"
	dnsRule := "unknown"
	cacheDomain, cachedRule, ok := GetCachedDnsItem(ip)
	if ok {
		domain = cacheDomain
		dnsRule = string(cachedRule.GetPolicy())
	}
	log.Infoln(log.FormatLog(log.RulePrefix, "ip:%v, rule:%v; domain:%v, rule:%v"), ip, result.GetPolicy(), domain, dnsRule)
	return result, nil
}

func NewTun(
	boostDns []string,
	remoteDns []string,
	localDns []string,
	defaultInterfaceName string,
	disableCache bool,
) (Config, error) {
	if len(boostDns) == 0 || len(remoteDns) == 0 || len(localDns) == 0 {
		return Config{}, fmt.Errorf("dns can't be empty")
	}
	ResetCache()
	dns, err := NewDnsDistribution(boostDns, remoteDns, localDns, defaultInterfaceName, disableCache)
	if err != nil {
		return Config{}, err
	}
	return Config{
		dns,
	}, nil
}
