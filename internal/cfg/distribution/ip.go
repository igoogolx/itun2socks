package distribution

import (
	"fmt"
	"github.com/igoogolx/itun2socks/internal/configuration"
	"github.com/igoogolx/itun2socks/internal/constants"
	geo2 "github.com/igoogolx/itun2socks/pkg/geo"
	"github.com/igoogolx/itun2socks/pkg/list"
)

type IpMatcher struct {
	Proxy  list.Lister
	Bypass list.Lister
}

type IpDistribution struct {
	Subnet       IpMatcher
	GeoIps       IpMatcher
	GeoSites     IpMatcher
	DefaultProxy bool
}

func NewIpDistribution(
	config configuration.IpItem,
) (IpDistribution, error) {
	d := IpDistribution{}
	d.Subnet = IpMatcher{
		Proxy: list.Lister{
			Items:  config.Subnet.Proxy,
			Mather: IsSubnetContainsIp,
		},
		Bypass: list.Lister{
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
		Proxy: list.Lister{
			Items:  proxyGeoIps,
			Mather: IsSubnetContainsIp,
		},
		Bypass: list.Lister{
			Items:  bypassGeoIps,
			Mather: IsSubnetContainsIp,
		},
	}

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