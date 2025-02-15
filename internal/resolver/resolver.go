package resolver

import (
	cResolver "github.com/Dreamacro/clash/component/resolver"
	_ "github.com/Dreamacro/clash/config"
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
	nameResolvers, err := parseNameServer(servers)
	if err != nil {
		return nil, err
	}
	for index, nameResolver := range nameResolvers {
		if nameResolver.Net == "dhcp" && nameResolver.Addr == "auto" {
			nameResolvers[index] = dns.NameServer{
				Net:       "dhcp",
				Addr:      defaultInterfaceName,
				Interface: defaultInterfaceName,
			}
		}
		if nameResolver.Addr == "system" {
			nameResolvers[index] = dns.NameServer{
				Net:       "system",
				Interface: defaultInterfaceName,
				Addr:      defaultInterfaceName,
			}
		}
	}
	return nameResolvers, err
}

//go:linkname parseNameServer github.com/Dreamacro/clash/config.parseNameServer
func parseNameServer(servers []string) ([]dns.NameServer, error)
