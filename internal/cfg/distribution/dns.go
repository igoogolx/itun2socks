package distribution

import (
	"github.com/igoogolx/itun2socks/internal/conn"
	"github.com/igoogolx/itun2socks/internal/resolver"
	"github.com/metacubex/mihomo/component/fakeip"
	metaC "github.com/metacubex/mihomo/constant"
	"github.com/metacubex/mihomo/dns"
)

func NewDnsDistribution(
	bootDns []string,
	remoteDns []string,
	localDns []string,
	disableCache bool,
	fakeIpPool *fakeip.Pool,
	defaultInterfaceName string,

) (DnsDistribution, error) {

	var err error
	dd := DnsDistribution{}

	//Boost
	boostDnsClient, err := resolver.New(bootDns, conn.DirectAdapter, defaultInterfaceName, disableCache)
	if err != nil {
		return DnsDistribution{}, err
	}
	boostDnsService := dns.NewService(boostDnsClient, dns.NewEnhancer(dns.EnhancerConfig{}))
	dd.Boost = SubDnsDistribution{
		Client:    boostDnsService,
		Addresses: bootDns,
		Resolvers: &boostDnsClient,
	}

	//Local
	localDnsClient, err := resolver.New(localDns, conn.DirectAdapter, defaultInterfaceName, disableCache)

	if err != nil {
		return DnsDistribution{}, err
	}

	localDnsService := dns.NewService(localDnsClient, dns.NewEnhancer(dns.EnhancerConfig{}))
	dd.Local = SubDnsDistribution{
		Addresses: localDns,
		Client:    localDnsService,
		Resolvers: &localDnsClient,
	}

	//Remote
	remoteDnsClient, err := resolver.New(remoteDns, conn.ProxyAdapter, defaultInterfaceName, disableCache)
	if err != nil {
		return DnsDistribution{}, err
	}

	remoteDnsEnhancerConfig := dns.EnhancerConfig{}

	if fakeIpPool != nil {

		remoteDnsEnhancerConfig = dns.EnhancerConfig{
			EnhancedMode:  metaC.DNSFakeIP,
			FakeIPPool:    fakeIpPool,
			FakeIPSkipper: &fakeip.Skipper{Mode: metaC.FilterBlackList},
		}

	}

	remoteDnsService := dns.NewService(remoteDnsClient, dns.NewEnhancer(remoteDnsEnhancerConfig))
	dd.Remote = SubDnsDistribution{
		Client:    remoteDnsService,
		Addresses: remoteDns,
		Resolvers: &remoteDnsClient,
	}

	return dd, nil
}

type SubDnsDistribution struct {
	Resolvers *dns.Resolvers
	Addresses []string
	Client    *dns.Service
}

func (s SubDnsDistribution) GetServers() []string {
	return s.Addresses
}

type DnsDistribution struct {
	Local  SubDnsDistribution
	Remote SubDnsDistribution
	Boost  SubDnsDistribution
}
