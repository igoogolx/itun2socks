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
) (DnsDistribution, error) {
	var err error
	var bootDnsServers []string
	for _, server := range bootDns {
		bootDnsServers = append(bootDnsServers, server+"#"+defaultInterfaceName)
	}
	boostDnsClient, err := resolver.New(bootDnsServers, defaultInterfaceName, func() (C.Proxy, error) {
		return conn.GetProxy(constants.RuleBypass)
	})
	if err != nil {
		return DnsDistribution{}, err
	}
	var localDnsServers []string
	for _, server := range bootDns {
		localDnsServers = append(localDns, server+"#"+defaultInterfaceName)
	}
	localDnsClient, err := resolver.New(localDnsServers, defaultInterfaceName, func() (C.Proxy, error) {
		return conn.GetProxy(constants.RuleBypass)
	})
	if err != nil {
		return DnsDistribution{}, err
	}
	dd := DnsDistribution{}
	dd.Local = SubDnsDistribution{
		Addresses: localDnsServers,
		Client:    localDnsClient,
	}

	remoteDnsClient, err := resolver.New(remoteDns, defaultInterfaceName, func() (C.Proxy, error) {
		return conn.GetProxy(constants.RuleProxy)
	})
	if err != nil {
		return DnsDistribution{}, err
	}
	dd.Remote = SubDnsDistribution{
		Client:    remoteDnsClient,
		Addresses: remoteDns,
	}

	dd.Boost = SubDnsDistribution{
		Client:    boostDnsClient,
		Addresses: bootDnsServers,
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
