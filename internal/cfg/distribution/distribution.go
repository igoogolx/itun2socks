package distribution

import (
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"github.com/igoogolx/itun2socks/internal/cfg/distribution/ruleEngine"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/pkg/log"
	"strings"
)

var dnsDomainCache, _ = lru.New(4 * 1024)
var dnsRuleCache, _ = lru.New(4 * 1024)

func GetCachedDnsItem(ip string) (string, constants.DnsType, bool) {
	rawCachedDomain, ok := dnsDomainCache.Get(ip)
	if !ok {
		return "", constants.DistributionLocalDns, false
	}
	rawCachedRule, ok := dnsRuleCache.Get(ip)
	if !ok {
		return "", constants.DistributionLocalDns, false
	}
	return rawCachedDomain.(string), constants.DnsType(rawCachedRule.(string)), true
}

func AddCachedDnsItem(ip, domain string, rule constants.DnsType) {
	dnsDomainCache.Add(ip, domain)
	dnsRuleCache.Add(ip, rule)
}

type Config struct {
	Dns        DnsDistribution
	RuleEngine *ruleEngine.Engine
	dnsTable   Cache
}

type Cache interface {
	Get(key interface{}) (interface{}, bool)
	Add(key interface{}, val interface{}) bool
}

func (c Config) GetDnsServerRule(ip string) (constants.RuleType, error) {
	if strings.Contains(c.Dns.Remote.Address, ip) {
		return constants.DistributionProxy, nil
	}
	return constants.DistributionBypass, fmt.Errorf("not found")
}

func (c Config) GetDnsRule(ip string) (constants.RuleType, error) {
	result := constants.DistributionBypass
	_, cachedRule, ok := GetCachedDnsItem(ip)
	if ok {
		if cachedRule == constants.DistributionLocalDns {
			result = constants.DistributionBypass
		} else {
			result = constants.DistributionProxy
		}
		return result, nil
	}
	return constants.DistributionBypass, fmt.Errorf("not found")
}

func (c Config) GetRule(ip string) constants.RuleType {

	result := constants.DistributionProxy

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

	//dns server
	result, err := c.GetDnsServerRule(ip)
	if err == nil {
		return result
	}

	//dns result
	result, err = c.GetDnsRule(ip)
	if err == nil {
		return result
	}

	rule, err := c.RuleEngine.Match(ip)
	if err == nil {
		if constants.RuleType(rule.Policy()) == constants.DistributionBypass {
			return constants.DistributionBypass
		} else if constants.RuleType(rule.Policy()) == constants.DistributionProxy {
			return constants.DistributionProxy
		}
	}

	return constants.DistributionProxy

}

func (c Config) GetDns(domain string) SubDnsDistribution {
	result := constants.DistributionRemoteDns
	if strings.Contains(c.Dns.Remote.Address, domain) {
		result = constants.DistributionBoostDns
	} else {
		var rule, err = c.RuleEngine.Match(domain)
		if err == nil {
			if constants.RuleType(rule.Policy()) == constants.DistributionBypass {
				result = constants.DistributionLocalDns
			} else if constants.RuleType(rule.Policy()) == constants.DistributionProxy {
				result = constants.DistributionRemoteDns
			}
		}
	}

	switch result {
	case constants.DistributionLocalDns:
		return c.Dns.Local
	case constants.DistributionRemoteDns:
		return c.Dns.Remote
	default:
		return c.Dns.Boost
	}
}

func New(
	boostDns string,
	remoteDns string,
	localDns string,
	ruleId string,
	rules []string,
	tunInterfaceName string,
	defaultInterfaceName string,
) (Config, error) {
	rEngine, err := ruleEngine.New(ruleId, rules)
	if err != nil {
		return Config{}, err
	}
	dns, err := NewDnsDistribution(boostDns, remoteDns, localDns, tunInterfaceName, defaultInterfaceName)
	if err != nil {
		return Config{}, err
	}
	return Config{
		Dns: dns, RuleEngine: rEngine,
	}, nil
}
