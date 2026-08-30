package distribution

import (
	"github.com/igoogolx/itun2socks/internal/conn"
	"github.com/igoogolx/itun2socks/internal/resolver"
	"github.com/metacubex/mihomo/component/fakeip"
	metaResolver "github.com/metacubex/mihomo/component/resolver"
	metaC "github.com/metacubex/mihomo/constant"
	"github.com/metacubex/mihomo/dns"
)

func NewDnsDistribution(
	bootDns []string,
	remoteDns []string,
	localDns []string,
	disableCache bool,
	fakeIpPool *fakeip.Pool,
) (DnsDistribution, error) {

	var err error
	dd := DnsDistribution{}

	//Boost
	boostDnsClient, err := resolver.New(bootDns, conn.DirectAdapter, disableCache)
	if err != nil {
		return DnsDistribution{}, err
	}
	boostDnsService := dns.NewService(boostDnsClient, dns.NewEnhancer(dns.EnhancerConfig{}))
	dd.Boost = SubDnsDistribution{
		Client:    boostDnsService,
		Addresses: bootDns,
	}

	//Local
	localDnsClient, err := resolver.New(localDns, conn.DirectAdapter, disableCache)

	if err != nil {
		return DnsDistribution{}, err
	}

	localDnsService := dns.NewService(localDnsClient, dns.NewEnhancer(dns.EnhancerConfig{}))
	dd.Local = SubDnsDistribution{
		Addresses: localDns,
		Client:    localDnsService,
	}

	//Remote
	remoteDnsClient, err := resolver.New(remoteDns, conn.ProxyAdapter, disableCache)
	if err != nil {
		return DnsDistribution{}, err
	}

	remoteDnsService := dns.NewService(remoteDnsClient, dns.NewEnhancer(dns.EnhancerConfig{
		IPv6:          false,
		EnhancedMode:  metaC.DNSFakeIP,
		FakeIPPool:    fakeIpPool,
		FakeIPSkipper: &fakeip.Skipper{Mode: metaC.FilterBlackList},
	}))
	dd.Remote = SubDnsDistribution{
		Client:    remoteDnsService,
		Addresses: remoteDns,
	}

	metaResolver.DefaultResolver = boostDnsClient
	return dd, nil
}

type SubDnsDistribution struct {
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
