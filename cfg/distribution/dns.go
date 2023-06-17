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
	bootDns string,
	remoteDns string,
	localDns string,
	config configuration.DnsItem,
) (DnsDistribution, error) {
	localAddress := localDns
	localDnsClient, err := resolver.NewClient(localAddress, bootDns, func(dohRemoteIp string) {
		AddCachedDnsItem(dohRemoteIp, localAddress, constants.DistributionLocalDns)
	})
	if err != nil {
		return DnsDistribution{}, err
	}
	dd := DnsDistribution{}
	localGeoSites, err := geo.LoadGeoSites(config.GeoSites.Local)
	if err != nil {
		return DnsDistribution{}, err
	}
	dd.Local = SubDnsDistribution{
		Address: localAddress,
		Client:  localDnsClient,
		Domains: list.New(
			config.Domains.Local,
			IsDomainMatchRule,
		),
		GeoSites: list.New(
			localGeoSites,
			IsContainsDomain,
		),
	}
	remoteGeoSites, err := geo.LoadGeoSites(config.GeoSites.Remote)
	if err != nil {
		return DnsDistribution{}, err
	}
	remoteDnsClient, err := resolver.NewClient(remoteDns, bootDns, func(dohRemoteIp string) {
		AddCachedDnsItem(dohRemoteIp, remoteDns, constants.DistributionRemoteDns)
	})
	if err != nil {
		return DnsDistribution{}, err
	}
	dd.Remote = SubDnsDistribution{
		Client:  remoteDnsClient,
		Address: remoteDns,
		Domains: list.New(
			config.Domains.Remote,
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
	dd.BoostNameserver = bootDns
	return dd, nil
}

type SubDnsDistribution struct {
	Domains  MatcherList
	GeoSites MatcherList
	Address  string
	Client   resolver.Client
}

type DnsDistribution struct {
	Local           SubDnsDistribution
	Remote          SubDnsDistribution
	BoostNameserver string
	Cache           Cache
}
