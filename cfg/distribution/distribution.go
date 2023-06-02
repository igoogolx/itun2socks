package distribution

import (
	"github.com/Dreamacro/clash/log"
	"github.com/igoogolx/itun2socks/components/resolver"
	"github.com/igoogolx/itun2socks/configuration"
	"github.com/igoogolx/itun2socks/constants"
	"golang.org/x/exp/slices"
	"strings"
)

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
	domain, ok := c.dnsTable.Get(ip)
	if slices.Contains(c.Dns.Local.Client.Nameservers(), ip) {
		result = constants.DistributionBypass
	} else if slices.Contains(c.Dns.Remote.Client.Nameservers(), ip) {
		result = constants.DistributionProxy
	} else if strings.Contains(c.TrueProxyServer, ip) {
		result = constants.DistributionBypass
	} else if ok && strings.Contains(c.TrueProxyServer, domain.(string)) {
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
					switch c.Ip.GeoSites.LookUp(domain.(string)) {
					case constants.DistributionProxy:
						result = constants.DistributionProxy
						break
					case constants.DistributionBypass:
						result = constants.DistributionBypass
						break
					case constants.DistributionNotFound:
						if c.Ip.DefaultProxy {
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

func (c Config) GetDns(domain string, isPrimary bool) resolver.Client {
	if isPrimary {
		return c.Dns.Local.Client
	}
	result := constants.DistributionPrimaryDns
	cacheResult, ok := c.Dns.Cache.Get(domain)
	if ok {
		result = cacheResult.(int)
	} else if IsDomainsContain(c.Dns.Remote.Client.Nameservers(), domain) {
		result = constants.DistributionPrimaryDns
	} else if strings.Contains(c.TrueProxyServer, domain) {
		result = constants.DistributionPrimaryDns
	} else {
		if c.Dns.Local.Domains.Has(domain) {
			result = constants.DistributionPrimaryDns
			log.Debugln("[Matching domain]: %v is from local domains", domain)
		} else if c.Dns.Local.GeoSites.Has(domain) {
			result = constants.DistributionPrimaryDns
			log.Debugln("[Matching domain]: %v is from local geo sites", domain)
		} else if c.Dns.Remote.Domains.Has(domain) {
			result = constants.DistributionSecondaryDns
			log.Debugln("[Matching domain]: %v is from remote domains", domain)
		} else if c.Dns.Remote.GeoSites.Has(domain) {
			result = constants.DistributionSecondaryDns
			log.Debugln("[Matching domain]: %v is from remote geo sites", domain)
		}
	}
	if result == constants.DistributionPrimaryDns {
		return c.Dns.Local.Client
	}
	return c.Dns.Remote.Client
}

func New(
	remoteDns string,
	localDns string,
	rule configuration.RuleCfg,
	trueProxyServer string,
	dnsTable Cache,
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
		Dns: dns, Ip: ip, TrueProxyServer: trueProxyServer, dnsTable: dnsTable,
	}, nil
}
