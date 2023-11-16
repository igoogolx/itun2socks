package distribution

import (
	cResolver "github.com/Dreamacro/clash/component/resolver"
	"github.com/igoogolx/itun2socks/internal/resolver"
)

/*

TODO: dhcp may be better as boost dns

boost dns works when local or remote is a doh server

local dns -> boost dns -> default network interface, will not be tracked
remote dns -> tun interface name -> boot dns -> proxy -> default network interface, will be tracked
boost dns -> default network interface, will not be tracked

*/

func NewDnsDistribution(
	bootDns string,
	remoteDns string,
	localDns string,
	tunInterfaceName string,
	defaultInterfaceName string,
) (DnsDistribution, error) {
	var err error
	bootDns = bootDns + "#" + defaultInterfaceName
	boostDnsClient, err := resolver.New([]string{bootDns}, []string{}, defaultInterfaceName)
	if err != nil {
		return DnsDistribution{}, err
	}
	localDns = localDns + "#" + defaultInterfaceName
	localDnsClient, err := resolver.New([]string{localDns}, []string{bootDns}, defaultInterfaceName)
	if err != nil {
		return DnsDistribution{}, err
	}
	dd := DnsDistribution{}
	dd.Local = SubDnsDistribution{
		Address: localDns,
		Client:  localDnsClient,
	}
	remoteDns = remoteDns + "#" + tunInterfaceName
	remoteDnsClient, err := resolver.New([]string{remoteDns}, []string{"udp://8.8.8.8#" + tunInterfaceName}, defaultInterfaceName)
	if err != nil {
		return DnsDistribution{}, err
	}
	dd.Remote = SubDnsDistribution{
		Client:  remoteDnsClient,
		Address: remoteDns,
	}

	dd.Boost = SubDnsDistribution{
		Client:  boostDnsClient,
		Address: bootDns,
	}

	cResolver.DefaultResolver = boostDnsClient
	return dd, nil
}

type SubDnsDistribution struct {
	Address string
	Client  cResolver.Resolver
}

type DnsDistribution struct {
	Local  SubDnsDistribution
	Remote SubDnsDistribution
	Boost  SubDnsDistribution
}
