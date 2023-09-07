package distribution

import (
	"github.com/Dreamacro/clash/component/resolver"
	"github.com/Dreamacro/clash/dns"
	"github.com/igoogolx/itun2socks/internal/configuration"
	"github.com/igoogolx/itun2socks/pkg/geo"
	"github.com/igoogolx/itun2socks/pkg/list"
	"strings"
)

func NewDnsDistribution(
	bootDns string,
	remoteDns string,
	localDns string,
	config configuration.DnsItem,
	tunDeviceName string,
) (DnsDistribution, error) {
	boostNameResolver := []dns.NameServer{{
		Net:  "tcp",
		Addr: bootDns,
	}}
	boostDnsClient := dns.NewResolver(dns.Config{
		Main: boostNameResolver,
	})
	localDnsNet := "tcp"
	if strings.Contains(remoteDns, "https") {
		localDnsNet = "https"
	}
	localDnsClient := dns.NewResolver(dns.Config{
		Main: []dns.NameServer{{
			Net:  localDnsNet,
			Addr: localDns,
		}},
		Default: boostNameResolver,
	})
	var err error
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
	remoteDnsNet := "tcp"
	if strings.Contains(remoteDns, "https") {
		remoteDnsNet = "https"
	}
	remoteDnsClient := dns.NewResolver(dns.Config{
		Main: []dns.NameServer{{
			Net:       remoteDnsNet,
			Addr:      remoteDns,
			Interface: tunDeviceName,
		}},
		Default: boostNameResolver,
	})
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

	dd.BoostNameserver = bootDns

	resolver.DefaultResolver = boostDnsClient
	return dd, nil
}

type SubDnsDistribution struct {
	Domains  MatcherList
	GeoSites MatcherList
	Address  string
	Client   resolver.Resolver
}

type DnsDistribution struct {
	Local           SubDnsDistribution
	Remote          SubDnsDistribution
	BoostNameserver string
}
