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

func (c Config) GetRule(ip string) (result constants.IpRule) {
	cacheResult, ok := c.Ip.Cache.Get(ip)
	if ok {
		return cacheResult.(constants.IpRule)
	}
	result = constants.DistributionBypass
	defer func() {
		c.Ip.Cache.Add(ip, result)
	}()
	//TODO: determine the type of true proxy server: ip or domain
	cachedDns, ok := GetCachedDnsItem(ip)
	if slices.Contains(c.Dns.Local.Client.Nameservers(), ip) {
		result = constants.DistributionBypass
	} else if slices.Contains(c.Dns.Remote.Client.Nameservers(), ip) {
		result = constants.DistributionProxy
	} else if strings.Contains(c.TrueProxyServer, ip) {
		result = constants.DistributionBypass
	} else if ok && strings.Contains(c.TrueProxyServer, cachedDns.Domain) {
		result = constants.DistributionBypass
	} else {
		switch c.Ip.Subnet.LookUp(ip) {
		case constants.DistributionProxy:
			result = constants.DistributionProxy
			break
		case constants.DistributionBypass:
			result = constants.DistributionBypass
			break
		case constants.DistributionNotFound:
			switch c.Ip.GeoIps.LookUp(ip) {
			case constants.DistributionProxy:
				result = constants.DistributionProxy
				break
			case constants.DistributionBypass:
				result = constants.DistributionBypass
				break
			case constants.DistributionNotFound:
				if ok {
					cacheItem, ok := GetCachedDnsItem(ip)
					if ok {
						if cacheItem.Rule == constants.DistributionLocalDns {
							result = constants.DistributionBypass
						} else {
							result = constants.DistributionProxy
						}
					}
				} else {
					if c.Ip.DefaultProxy {
						result = constants.DistributionProxy
					}
				}
			}
		}
	}

	return result
}

func (c Config) GetDns(domain string, isLocal bool) (resolver.Client, constants.DnsRule) {
	if isLocal {
		return c.Dns.Local.Client, constants.DistributionLocalDns
	}
	result := constants.DistributionLocalDns
	cacheResult, ok := c.Dns.Cache.Get(domain)
	if ok {
		result = cacheResult.(constants.DnsRule)
	} else if IsDomainsContain(c.Dns.Remote.Client.Nameservers(), domain) {
		result = constants.DistributionLocalDns
	} else if strings.Contains(c.TrueProxyServer, domain) {
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
	c.Dns.Cache.Add(domain, result)
	if result == constants.DistributionLocalDns {
		return c.Dns.Local.Client, constants.DistributionLocalDns
	}
	return c.Dns.Remote.Client, constants.DistributionRemoteDns
}

func New(
	remoteDns string,
	localDns string,
	rule configuration.RuleCfg,
	trueProxyServer string,
) (Config, error) {
	dns, err := NewDnsDistribution(remoteDns, localDns, rule.Dns)
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
