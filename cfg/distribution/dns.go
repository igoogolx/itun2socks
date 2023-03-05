package distribution

import (
	"fmt"
	"github.com/hashicorp/golang-lru"
	"github.com/igoogolx/itun2socks/components/geo"
	"github.com/igoogolx/itun2socks/components/list"
	"github.com/igoogolx/itun2socks/components/resolver"
	"github.com/igoogolx/itun2socks/configuration"
	"github.com/igoogolx/itun2socks/constants"
)

func NewDnsDistribution(
	config configuration.DnsItem,
	geoDataDir string,
) (DnsDistribution, error) {
	localAddress := config.Local.Address
	localDnsClient, err := resolver.NewClient(localAddress)
	if err != nil {
		return DnsDistribution{}, err
	}
	dd := DnsDistribution{}
	localGeoSites, err := geo.LoadGeoSites(geoDataDir, config.Local.GeoSites)
	if err != nil {
		return DnsDistribution{}, err
	}
	dd.Local = SubDnsDistribution{
		Address: localAddress,
		Client:  localDnsClient,
		Domains: list.New(
			config.Local.Domains,
			IsDomainMatchRule,
		),
		GeoSites: list.New(
			localGeoSites,
			IsContainsDomain,
		),
	}
	remoteGeoSites, err := geo.LoadGeoSites(geoDataDir, config.Remote.GeoSites)
	if err != nil {
		return DnsDistribution{}, err
	}
	remoteDnsClient, err := resolver.NewClient(config.Remote.Address)
	if err != nil {
		return DnsDistribution{}, err
	}
	dd.Remote = SubDnsDistribution{
		Client:  remoteDnsClient,
		Address: config.Remote.Address,
		Domains: list.New(
			config.Remote.Domains,
			IsDomainMatchRule,
		),
		GeoSites: list.New(
			remoteGeoSites,
			IsContainsDomain,
		),
	}
	dd.Cache, err = lru.New(constants.CacheSize)
	if err != nil {
		return DnsDistribution{}, fmt.Errorf("fail to init dns cache,err:%v", err)
	}
	return dd, nil
}

type SubDnsDistribution struct {
	Domains  MatcherList
	GeoSites MatcherList
	Address  string
	Client   resolver.Client
}

type DnsDistribution struct {
	Local  SubDnsDistribution
	Remote SubDnsDistribution
	Cache  Cache
}
