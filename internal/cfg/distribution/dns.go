package distribution

import (
	"github.com/Dreamacro/clash/component/fakeip"
	cResolver "github.com/Dreamacro/clash/component/resolver"
	C "github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/conn"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/internal/resolver"
)

func NewDnsDistribution(
	bootDns []string,
	remoteDns []string,
	localDns []string,
	defaultInterfaceName string,
	disableCache bool,
	fakeIpPool *fakeip.Pool,
) (DnsDistribution, error) {

	var err error
	dd := DnsDistribution{}

	//Boost
	boostDnsClient, err := resolver.New(bootDns, defaultInterfaceName, func() (C.Proxy, error) {
		return conn.GetProxy(constants.PolicyDirect)
	}, disableCache, nil)
	if err != nil {
		return DnsDistribution{}, err
	}
	dd.Boost = SubDnsDistribution{
		Client:    boostDnsClient,
		Addresses: bootDns,
	}

	//Local
	localDnsClient, err := resolver.New(localDns, defaultInterfaceName, func() (C.Proxy, error) {
		return conn.GetProxy(constants.PolicyDirect)
	}, disableCache, nil)
	if err != nil {
		return DnsDistribution{}, err
	}
	dd.Local = SubDnsDistribution{
		Addresses: localDns,
		Client:    localDnsClient,
	}

	//Remote
	remoteDnsClient, err := resolver.New(remoteDns, defaultInterfaceName, func() (C.Proxy, error) {
		return conn.GetProxy(constants.PolicyProxy)
	}, disableCache, fakeIpPool)
	if err != nil {
		return DnsDistribution{}, err
	}
	dd.Remote = SubDnsDistribution{
		Client:    remoteDnsClient,
		Addresses: remoteDns,
	}

	cResolver.DefaultResolver = boostDnsClient
	return dd, nil
}

type SubDnsDistribution struct {
	Addresses []string
	Client    cResolver.Resolver
}

func (s SubDnsDistribution) GetServers() []string {
	return s.Client.GetServers()
}

type DnsDistribution struct {
	Local  SubDnsDistribution
	Remote SubDnsDistribution
	Boost  SubDnsDistribution
}
