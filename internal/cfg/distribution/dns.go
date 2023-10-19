package distribution

import (
	cResolver "github.com/Dreamacro/clash/component/resolver"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/internal/resolver"
)

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
		Type:    constants.DistributionLocalDns,
	}
	remoteDns = remoteDns + "#" + tunInterfaceName
	remoteDnsClient, err := resolver.New([]string{remoteDns}, []string{"udp://8.8.8.8#" + tunInterfaceName}, defaultInterfaceName)
	if err != nil {
		return DnsDistribution{}, err
	}
	dd.Remote = SubDnsDistribution{
		Client:  remoteDnsClient,
		Address: remoteDns,
		Type:    constants.DistributionRemoteDns,
	}

	dd.Boost = SubDnsDistribution{
		Client:  boostDnsClient,
		Address: bootDns,
		Type:    constants.DistributionBoostDns,
	}

	cResolver.DefaultResolver = boostDnsClient
	return dd, nil
}

type SubDnsDistribution struct {
	Address string
	Type    constants.DnsType
	Client  cResolver.Resolver
}

type DnsDistribution struct {
	Local  SubDnsDistribution
	Remote SubDnsDistribution
	Boost  SubDnsDistribution
}
