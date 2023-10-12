package distribution

import (
	"github.com/Dreamacro/clash/component/resolver"
	lru "github.com/hashicorp/golang-lru"
	rule2 "github.com/igoogolx/itun2socks/internal/cfg/distribution/rule"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/pkg/log"
	"strings"
)

var dnsCache, _ = lru.New(1000)

func getRuleStr(rule constants.RuleType) string {
	if rule == constants.DistributionBypass {
		return "direct"
	}
	return "proxy"
}

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
	Dns      DnsDistribution
	Rules    []rule2.Rule
	dnsTable Cache
}

type Cache interface {
	Get(key interface{}) (interface{}, bool)
	Add(key interface{}, val interface{}) bool
}

func (c Config) GetDnsServerRule(ip string) constants.RuleType {
	result := constants.DistributionNotFound
	if strings.Contains(c.Dns.Remote.Address, ip) {
		result = constants.DistributionProxy
	}
	return result
}

func (c Config) GetDnsRule(ip string) constants.RuleType {
	result := constants.DistributionNotFound
	cacheItem, ok := GetCachedDnsItem(ip)
	if ok {
		if cacheItem.Rule == constants.DistributionLocalDns {
			result = constants.DistributionBypass
		} else {
			result = constants.DistributionProxy
		}
	}
	return result
}

func (c Config) GetRule(ip string) constants.RuleType {

	rule := constants.DistributionBypass

	defer func(latestIp string) {
		domain := "unknown"
		dnsRule := "unknown"
		cacheItem, ok := GetCachedDnsItem(latestIp)
		if ok {
			domain = cacheItem.Domain
			dnsRule = string(cacheItem.Rule)
		}
		log.Infoln(log.FormatLog(log.RulePrefix, "ip:%v, rule:%v; domain:%v, rule:%v"), latestIp, getRuleStr(rule), domain, dnsRule)
	}(ip)

	//dns server
	rule = c.GetDnsServerRule(ip)
	if rule != constants.DistributionNotFound {
		return rule
	}

	//dns rule
	rule = c.GetDnsRule(ip)
	if rule != constants.DistributionNotFound {
		return rule
	}

	for _, rule := range c.Rules {
		if rule.Match(ip) {
			if rule.Policy() == "bypass" {
				return constants.DistributionBypass
			}
			if rule.Policy() == "proxy" {
				return constants.DistributionProxy
			}
		}
	}

	return constants.DistributionProxy

}

func (c Config) GetDns(domain string) (resolver.Resolver, constants.DnsType) {
	result := constants.DistributionLocalDns
	if strings.Contains(c.Dns.Remote.Address, domain) {
		result = constants.DistributionBoostDns
	} else {
		for _, rule := range c.Rules {
			if rule.Match(domain) {
				if rule.Policy() == "bypass" {
					result = constants.DistributionLocalDns
				}
				if rule.Policy() == "proxy" {
					result = constants.DistributionRemoteDns
				}
			}
		}

	}

	if result == constants.DistributionLocalDns {
		return c.Dns.Local.Client, constants.DistributionLocalDns
	}
	if result == constants.DistributionBoostDns {
		return c.Dns.Boost.Client, constants.DistributionBoostDns
	}
	return c.Dns.Remote.Client, constants.DistributionRemoteDns
}

func New(
	boostDns string,
	remoteDns string,
	localDns string,
	rule string,
	tunInterfaceName string,
	defaultInterfaceName string,
) (Config, error) {
	rules, err := rule2.Parse(rule)
	if err != nil {
		return Config{}, err
	}
	dns, err := NewDnsDistribution(boostDns, remoteDns, localDns, tunInterfaceName, defaultInterfaceName)
	if err != nil {
		return Config{}, err
	}
	return Config{
		Dns: dns, dnsTable: dnsCache, Rules: rules,
	}, nil
}
