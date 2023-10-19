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
		return "", constants.LocalDns, false
	}
	rawCachedRule, ok := dnsRuleCache.Get(ip)
	if !ok {
		return "", constants.LocalDns, false
	}
	return rawCachedDomain.(string), rawCachedRule.(constants.DnsType), true
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
		return constants.RuleProxy, nil
	}
	return constants.RuleBypass, fmt.Errorf("not found")
}

func (c Config) GetDnsRule(ip string) (constants.RuleType, error) {
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

func (c Config) GetRule(ip string) constants.RuleType {

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
		if constants.RuleType(rule.Policy()) == constants.RuleBypass {
			return constants.RuleBypass
		} else if constants.RuleType(rule.Policy()) == constants.RuleProxy {
			return constants.RuleProxy
		}
	}

	return constants.RuleProxy

}

func (c Config) GetDns(domain string) SubDnsDistribution {
	result := constants.RemoteDns
	if strings.Contains(c.Dns.Remote.Address, domain) {
		result = constants.BoostDns
	} else {
		var rule, err = c.RuleEngine.Match(domain)
		if err == nil {
			if constants.RuleType(rule.Policy()) == constants.RuleBypass {
				result = constants.LocalDns
			} else if constants.RuleType(rule.Policy()) == constants.RuleProxy {
				result = constants.RemoteDns
			}
		}
	}

	switch result {
	case constants.LocalDns:
		return c.Dns.Local
	case constants.RemoteDns:
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
