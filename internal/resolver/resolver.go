package resolver

import (
	cResolver "github.com/Dreamacro/clash/component/resolver"
	"github.com/Dreamacro/clash/config"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/dns"
	_ "unsafe"
)

func New(mainServer []string, defaultInterfaceName string, getDialer func() (C.Proxy, error), disableCache bool) (cResolver.Resolver, error) {
	mainNameResolver, err := parse(mainServer, defaultInterfaceName)
	if err != nil {
		return nil, err
	}

	mainDnsClient := dns.NewResolver(dns.Config{
		Main:         mainNameResolver,
		GetDialer:    getDialer,
		DisableCache: disableCache,
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
