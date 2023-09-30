package distribution

import (
	cResolver "github.com/Dreamacro/clash/component/resolver"
	"github.com/igoogolx/itun2socks/internal/configuration"
	"github.com/igoogolx/itun2socks/internal/resolver"
	"github.com/igoogolx/itun2socks/pkg/geo"
	"github.com/igoogolx/itun2socks/pkg/list"
	"github.com/igoogolx/itun2socks/pkg/network_iface"
	"strings"
)

func NewDnsDistribution(
	bootDns string,
	remoteDns string,
	localDns string,
	config configuration.DnsItem,
	tunDeviceName string,
) (DnsDistribution, error) {
	var err error
	bootDns = bootDns + "#" + network_iface.GetDefaultInterfaceName()
	boostDnsClient, err := resolver.New([]string{bootDns}, []string{})
	if err != nil {
		return DnsDistribution{}, err
	}
	localDns = localDns + "#" + network_iface.GetDefaultInterfaceName()
	localDnsClient, err := resolver.New([]string{localDns}, []string{bootDns})
	if err != nil {
		return DnsDistribution{}, err
	}
	dd := DnsDistribution{}
	localGeoSites, err := geo.LoadGeoSites(config.GeoSites.Local)
	if err != nil {
		return DnsDistribution{}, err
	}
	dd.Local = SubDnsDistribution{
		Address: localDns,
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
	remoteDns = remoteDns + "#" + tunDeviceName
	remoteDnsClient, err := resolver.New([]string{remoteDns}, []string{"udp://8.8.8.8#" + tunDeviceName})
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

	dd.Boost = SubDnsDistribution{
		Client:  boostDnsClient,
		Address: bootDns,
		Domains: list.New(
			[]string{remoteDns},
			strings.Contains,
		),
		GeoSites: list.New(
			[]string{},
			IsContainsDomain,
		),
	}

	cResolver.DefaultResolver = boostDnsClient
	return dd, nil
}

type SubDnsDistribution struct {
	Domains  MatcherList
	GeoSites MatcherList
	Address  string
	Client   cResolver.Resolver
}

type DnsDistribution struct {
	Local  SubDnsDistribution
	Remote SubDnsDistribution
	Boost  SubDnsDistribution
}
