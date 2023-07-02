package distribution

import (
	"github.com/Dreamacro/clash/log"
	lru "github.com/hashicorp/golang-lru"
	"github.com/igoogolx/itun2socks/components/resolver"
	"github.com/igoogolx/itun2socks/configuration"
	"github.com/igoogolx/itun2socks/constants"
	"golang.org/x/exp/slices"
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
	Dns             DnsDistribution
	Ip              IpDistribution
	TrueProxyServer string
	dnsTable        Cache
}

type Cache interface {
	Get(key interface{}) (interface{}, bool)
	Add(key interface{}, val interface{}) bool
}

func (c Config) GetDnsServerRule(ip string) constants.IpRule {
	result := constants.DistributionNotFound
	if slices.Contains(c.Dns.Local.Client.Nameservers(), ip) {
		result = constants.DistributionBypass
	} else if strings.Contains(c.Dns.BoostNameserver, ip) {
		result = constants.DistributionBypass
	} else if IsDomainsContain(c.Dns.Remote.Client.Nameservers(), ip) {
		result = constants.DistributionProxy
	}
	return result
}

func (c Config) GetTrueProxyServerRule(ip string) constants.IpRule {
	result := constants.DistributionNotFound
	cachedDns, ok := GetCachedDnsItem(ip)
	if strings.Contains(c.TrueProxyServer, ip) || (ok && strings.Contains(c.TrueProxyServer, cachedDns.Domain)) {
		result = constants.DistributionBypass
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

	//dns server
	rule = c.GetDnsServerRule(ip)
	if rule != constants.DistributionNotFound {
		return rule
	}

	//true proxy server
	rule = c.GetTrueProxyServerRule(ip)
	if rule != constants.DistributionNotFound {
		return rule
	}

	//subnet rule
	rule = c.GetTrueProxyServerRule(ip)
	if rule != constants.DistributionNotFound {
		return rule
	}

	//geo rule
	rule = c.GetTrueProxyServerRule(ip)
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

func (c Config) GetDns(domain string, isLocal bool) (resolver.Client, constants.DnsRule) {
	result := constants.DistributionLocalDns
	if isLocal {
		log.Debugln("[Matching domain]: %v is from local", domain)
		result = constants.DistributionLocalDns
	} else if strings.Contains(c.TrueProxyServer, domain) {
		log.Debugln("[Matching domain]: %v is from true proxy server", domain)
		result = constants.DistributionLocalDns
	} else {
		if c.Dns.Local.Domains.Has(domain) {
			result = constants.DistributionLocalDns
			log.Debugln("[Matching domain]: %v is from local domains", domain)
		} else if c.Dns.Local.GeoSites.Has(domain) {
			result = constants.DistributionLocalDns
			log.Debugln("[Matching domain]: %v is from local geo sites", domain)
		} else if c.Dns.Remote.Domains.Has(domain) {
			result = constants.DistributionRemoteDns
			log.Debugln("[Matching domain]: %v is from remote domains", domain)
		} else if c.Dns.Remote.GeoSites.Has(domain) {
			result = constants.DistributionRemoteDns
			log.Debugln("[Matching domain]: %v is from remote geo sites", domain)
		}
	}
	if result == constants.DistributionLocalDns {
		return c.Dns.Local.Client, constants.DistributionLocalDns
	}
	return c.Dns.Remote.Client, constants.DistributionRemoteDns
}

func New(
	boostDns string,
	remoteDns string,
	localDns string,
	rule configuration.RuleCfg,
	trueProxyServer string,
) (Config, error) {
	dns, err := NewDnsDistribution(boostDns, remoteDns, localDns, rule.Dns)
	if err != nil {
		return Config{}, err
	}
	ip, err := NewIpDistribution(rule.Ip)
	if err != nil {
		return Config{}, err
	}
	return Config{
		Dns: dns, Ip: ip, TrueProxyServer: trueProxyServer, dnsTable: dnsCache,
	}, nil
}
