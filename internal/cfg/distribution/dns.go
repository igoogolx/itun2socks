package distribution

import (
	"github.com/Dreamacro/clash/component/resolver"
	"github.com/Dreamacro/clash/dns"
	"github.com/igoogolx/itun2socks/internal/configuration"
	"github.com/igoogolx/itun2socks/pkg/geo"
	"github.com/igoogolx/itun2socks/pkg/list"
)

func NewDnsDistribution(
	bootDns string,
	remoteDns string,
	localDns string,
	config configuration.DnsItem,
	tunDeviceName string,
) (DnsDistribution, error) {
	localDnsClient := dns.NewResolver(dns.Config{
		Main: []dns.NameServer{{
			Net:  "tcp",
			Addr: localDns,
		}},
		Default: []dns.NameServer{
			{
				Net:  "tcp",
				Addr: bootDns,
			},
		},
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
	remoteDnsClient := dns.NewResolver(dns.Config{
		Main: []dns.NameServer{{
			Net:       "tcp",
			Addr:      remoteDns,
			Interface: tunDeviceName,
		}},
		Default: []dns.NameServer{
			{
				Net:  "tcp",
				Addr: bootDns,
			},
		},
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
