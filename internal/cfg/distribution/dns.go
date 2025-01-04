package distribution

import (
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
) (DnsDistribution, error) {

	var err error
	dd := DnsDistribution{}

	//Boost
	var bootDnsServers []string
	for _, server := range bootDns {
		bootDnsServers = append(bootDnsServers, server+"#"+defaultInterfaceName)
	}
	boostDnsClient, err := resolver.New(bootDnsServers, defaultInterfaceName, func() (C.Proxy, error) {
		return conn.GetProxy(constants.PolicyDirect)
	}, disableCache)
	if err != nil {
		return DnsDistribution{}, err
	}
	dd.Boost = SubDnsDistribution{
		Client:    boostDnsClient,
		Addresses: bootDnsServers,
	}

	//Local
	var localDnsServers []string
	for _, server := range localDns {
		localDnsServers = append(localDnsServers, server+"#"+defaultInterfaceName)
	}
	localDnsClient, err := resolver.New(localDnsServers, defaultInterfaceName, func() (C.Proxy, error) {
		return conn.GetProxy(constants.PolicyDirect)
	}, disableCache)
	if err != nil {
		return DnsDistribution{}, err
	}
	dd.Local = SubDnsDistribution{
		Addresses: localDnsServers,
		Client:    localDnsClient,
	}

	//Remote
	remoteDnsClient, err := resolver.New(remoteDns, defaultInterfaceName, func() (C.Proxy, error) {
		return conn.GetProxy(constants.PolicyProxy)
	}, disableCache)
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

type DnsDistribution struct {
	Local  SubDnsDistribution
	Remote SubDnsDistribution
	Boost  SubDnsDistribution
}
