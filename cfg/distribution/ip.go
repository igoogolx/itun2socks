package distribution

import (
	"errors"
	"fmt"
	"github.com/hashicorp/golang-lru"
	geo2 "github.com/igoogolx/itun2socks/components/geo"
	"github.com/igoogolx/itun2socks/components/list"
	"github.com/igoogolx/itun2socks/configuration"
	"github.com/igoogolx/itun2socks/constants"
)

type IpMatcher struct {
	Proxy  list.List
	Bypass list.List
}

type IpDistribution struct {
	Subnet       IpMatcher
	GeoIps       IpMatcher
	GeoSites     IpMatcher
	Cache        Cache
	DefaultProxy bool
}

func NewIpDistribution(
	config configuration.IpItem,
) (IpDistribution, error) {
	d := IpDistribution{}
	d.Subnet = IpMatcher{
		Proxy: list.List{
			Items:  config.Subnet.Proxy,
			Mather: IsSubnetContainsIp,
		},
		Bypass: list.List{
			Items:  config.Subnet.Bypass,
			Mather: IsSubnetContainsIp,
		},
	}
	proxyGeoIps, err := geo2.LoadGeoIPs(config.GeoIps.Proxy)
	if err != nil {
		return IpDistribution{}, fmt.Errorf("fail to parse proxy ip geo ips :%v", err)
	}
	bypassGeoIps, err := geo2.LoadGeoIPs(config.GeoIps.Bypass)
	if err != nil {
		return IpDistribution{}, fmt.Errorf("fail to parse bypass ip geo ips :%v", err)
	}
	d.GeoIps = IpMatcher{
		Proxy: list.List{
			Items:  proxyGeoIps,
			Mather: IsSubnetContainsIp,
		},
		Bypass: list.List{
			Items:  bypassGeoIps,
			Mather: IsSubnetContainsIp,
		},
	}

	proxyGeoSites, err := geo2.LoadGeoSites(config.GeoSites.Proxy)
	if err != nil {
		return IpDistribution{}, fmt.Errorf("fail to parse proxy ip geo sites :%v", err)
	}
	bypassGeoSites, err := geo2.LoadGeoSites(config.GeoSites.Bypass)
	if err != nil {
		return IpDistribution{}, fmt.Errorf("fail to parse bypass ip geo sites :%v", err)
	}
	d.GeoSites = IpMatcher{
		Proxy: list.List{
			Items:  proxyGeoSites,
			Mather: IsContainsDomain,
		},
		Bypass: list.List{
			Items:  bypassGeoSites,
			Mather: IsContainsDomain,
		},
	}

	ipCache, err := lru.New(constants.CacheSize)
	if err != nil {
		return IpDistribution{}, errors.New("fail to init proxy ip region cache")
	}
	d.Cache = ipCache
	d.DefaultProxy = config.DefaultProxy
	return d, nil
}

func (s *IpMatcher) LookUp(ip string) constants.IpRule {
	if s.Proxy.Has(ip) {
		return constants.DistributionProxy
	}
	if s.Bypass.Has(ip) {
		return constants.DistributionBypass
	}
	return constants.DistributionNotFound
}
