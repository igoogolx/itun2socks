package distribution

import (
	cResolver "github.com/Dreamacro/clash/component/resolver"
	C "github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/conn"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/internal/resolver"
)

func NewDnsDistribution(
	bootDns string,
	remoteDns string,
	localDns string,
	defaultInterfaceName string,
) (DnsDistribution, error) {
	var err error
	bootDns = bootDns + "#" + defaultInterfaceName
	boostDnsClient, err := resolver.New([]string{bootDns}, defaultInterfaceName, func() C.Proxy {
		return conn.GetProxy(constants.RuleBypass)
	})
	if err != nil {
		return DnsDistribution{}, err
	}
	localDns = localDns + "#" + defaultInterfaceName
	localDnsClient, err := resolver.New([]string{localDns}, defaultInterfaceName, func() C.Proxy {
		return conn.GetProxy(constants.RuleBypass)
	})
	if err != nil {
		return DnsDistribution{}, err
	}
	dd := DnsDistribution{}
	dd.Local = SubDnsDistribution{
		Address: localDns,
		Client:  localDnsClient,
	}

	remoteDnsClient, err := resolver.New([]string{remoteDns}, defaultInterfaceName, func() C.Proxy {
		return conn.GetProxy(constants.RuleProxy)
	})
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
