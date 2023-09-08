package distribution

import (
	"github.com/Dreamacro/clash/component/resolver"
	lru "github.com/hashicorp/golang-lru"
	"github.com/igoogolx/itun2socks/internal/configuration"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/pkg/log"
	"strings"
)

var dnsCache, _ = lru.New(1000)

func getRuleStr(rule constants.IpRule) string {
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

func AddCachedDnsItem(ip, domain string, rule constants.DnsRule) {
	dnsCache.Add(ip, CacheItem{
		Domain: domain,
		Rule:   rule,
	})
}

type CacheItem struct {
	Domain string
	Rule   constants.DnsRule
}

type Config struct {
	Dns      DnsDistribution
	Ip       IpDistribution
	dnsTable Cache
}

type Cache interface {
	Get(key interface{}) (interface{}, bool)
	Add(key interface{}, val interface{}) bool
}

func (c Config) GetDnsServerRule(ip string) constants.IpRule {
	result := constants.DistributionNotFound
	if strings.Contains(c.Dns.Remote.Address, ip) {
		result = constants.DistributionProxy
	}
	return result
}

func (c Config) GetSubnetRule(ip string) constants.IpRule {
	return c.Ip.Subnet.LookUp(ip)
}

func (c Config) GetGeoRule(ip string) constants.IpRule {
	return c.Ip.GeoIps.LookUp(ip)
}

func (c Config) GetDnsRule(ip string) constants.IpRule {
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

func (c Config) GetRule(ip string) constants.IpRule {

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

	//subnet rule
	rule = c.GetSubnetRule(ip)
	if rule != constants.DistributionNotFound {
		return rule
	}

	//geo rule
	rule = c.GetGeoRule(ip)
	if rule != constants.DistributionNotFound {
		return rule
	}

	//dns rule
	rule = c.GetDnsRule(ip)
	if rule != constants.DistributionNotFound {
		return rule
	}

	//default rule
	if c.Ip.DefaultProxy {
		rule = constants.DistributionProxy
	}

	return rule

}

func (c Config) GetDns(domain string) (resolver.Resolver, constants.DnsRule) {
	result := constants.DistributionLocalDns
	if c.Dns.Boost.Domains.Has(domain) {
		result = constants.DistributionBoostDns
		log.Debugln(log.FormatLog(log.RulePrefix, "%v is from boost domains"), domain)
	} else if c.Dns.Local.Domains.Has(domain) {
		result = constants.DistributionLocalDns
		log.Debugln(log.FormatLog(log.RulePrefix, "%v is from local domains"), domain)
	} else if c.Dns.Local.GeoSites.Has(domain) {
		result = constants.DistributionLocalDns
		log.Debugln(log.FormatLog(log.RulePrefix, "%v is from local geo sites"), domain)
	} else if c.Dns.Remote.Domains.Has(domain) {
		result = constants.DistributionRemoteDns
		log.Debugln(log.FormatLog(log.RulePrefix, "%v is from remote domains"), domain)
	} else if c.Dns.Remote.GeoSites.Has(domain) {
		result = constants.DistributionRemoteDns
		log.Debugln(log.FormatLog(log.RulePrefix, "%v is from remote geo sites"), domain)
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
	rule configuration.RuleCfg,
	tunDeviceName string,
) (Config, error) {
	dns, err := NewDnsDistribution(boostDns, remoteDns, localDns, rule.Dns, tunDeviceName)
	if err != nil {
		return Config{}, err
	}
	ip, err := NewIpDistribution(rule.Ip)
	if err != nil {
		return Config{}, err
	}
	return Config{
		Dns: dns, Ip: ip, dnsTable: dnsCache,
	}, nil
}
