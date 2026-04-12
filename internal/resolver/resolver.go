package resolver

import (
	_ "unsafe"

	"github.com/igoogolx/itun2socks/pkg/clash/component/fakeip"
	cResolver "github.com/igoogolx/itun2socks/pkg/clash/component/resolver"
	"github.com/igoogolx/itun2socks/pkg/clash/config"
	C "github.com/igoogolx/itun2socks/pkg/clash/constant"
	"github.com/igoogolx/itun2socks/pkg/clash/dns"
)

func New(mainServer []string, defaultInterfaceName string, getDialer func() (C.Proxy, error), disableCache bool, fakeIpPool *fakeip.Pool) (cResolver.Resolver, error) {
	mainNameResolver, err := parse(mainServer, defaultInterfaceName)
	if err != nil {
		return nil, err
	}

	mainDnsClient := dns.NewResolver(dns.Config{
		Main:         mainNameResolver,
		GetDialer:    getDialer,
		DisableCache: disableCache,
		Pool:         fakeIpPool,
	})

	return mainDnsClient, nil
}

func parse(servers []string, defaultInterfaceName string) ([]dns.NameServer, error) {
	nameResolvers, err := config.ParseNameServer(servers)
	if err != nil {
		return nil, err
	}
	for index, nameResolver := range nameResolvers {
		//FIXME: remove dhcp
		if nameResolver.Net == "system" || nameResolver.Net == "dhcp" {
			nameResolvers[index] = dns.NameServer{
				Net:       "system",
				Interface: defaultInterfaceName,
				Addr:      defaultInterfaceName,
			}
		}
	}
	return nameResolvers, err
}
