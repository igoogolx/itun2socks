package distribution

import (
	"fmt"
	"github.com/Dreamacro/clash/component/resolver"
	lru "github.com/hashicorp/golang-lru"
	"github.com/igoogolx/itun2socks/internal/cfg/distribution/ruleEngine"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/pkg/log"
	"strings"
)

var dnsCache, _ = lru.New(1000)

func GetCachedDnsItem(ip string) (CacheItem, bool) {
	cacheItem, ok := dnsCache.Get(ip)
	if ok {
		cacheResult, ok := cacheItem.(CacheItem)
		return cacheResult, ok
	}
	return CacheItem{}, false
}

func AddCachedDnsItem(ip, domain string, rule constants.DnsType) {
	dnsCache.Add(ip, CacheItem{
		Domain: domain,
		Rule:   rule,
	})
}

type CacheItem struct {
	Domain string
	Rule   constants.DnsType
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
	cacheItem, ok := GetCachedDnsItem(ip)
	if ok {
		if cacheItem.Rule == constants.DistributionLocalDns {
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
		cacheItem, ok := GetCachedDnsItem(latestIp)
		if ok {
			domain = cacheItem.Domain
			dnsRule = string(cacheItem.Rule)
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

func (c Config) GetDns(domain string) (resolver.Resolver, constants.DnsType) {
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
		return c.Dns.Local.Client, constants.DistributionLocalDns
	case constants.DistributionRemoteDns:
		return c.Dns.Remote.Client, constants.DistributionRemoteDns
	default:
		return c.Dns.Boost.Client, constants.DistributionBoostDns

	}
}

func New(
	boostDns string,
	remoteDns string,
	localDns string,
	rule string,
	tunInterfaceName string,
	defaultInterfaceName string,
) (Config, error) {
	ruleEngine, err := ruleEngine.New(rule)
	if err != nil {
		return Config{}, err
	}
	dns, err := NewDnsDistribution(boostDns, remoteDns, localDns, tunInterfaceName, defaultInterfaceName)
	if err != nil {
		return Config{}, err
	}
	return Config{
		Dns: dns, dnsTable: dnsCache, RuleEngine: ruleEngine,
	}, nil
}
