package resolver

import (
	cResolver "github.com/Dreamacro/clash/component/resolver"
	"github.com/Dreamacro/clash/dns"
	"github.com/igoogolx/itun2socks/pkg/network_iface"
	_ "unsafe"
)

func New(mainServer, defaultServer []string) (cResolver.Resolver, error) {
	defaultNameResolver, err := parse(defaultServer)
	if err != nil {
		return nil, err
	}
	mainNameResolver, err := parse(mainServer)
	if err != nil {
		return nil, err
	}

	mainDnsClient := dns.NewResolver(dns.Config{
		Main:    mainNameResolver,
		Default: defaultNameResolver,
	})

	return mainDnsClient, nil
}

func parse(servers []string) ([]dns.NameServer, error) {
	nameResolvers, err := parseNameServer(servers)
	if err != nil {
		return nil, err
	}
	for index, nameResolver := range nameResolvers {
		if nameResolver.Net == "dhcp" && nameResolver.Addr == "auto" {
			nameResolvers[index] = dns.NameServer{
				Net:       "dhcp",
				Addr:      network_iface.GetDefaultInterfaceName(),
				Interface: network_iface.GetDefaultInterfaceName(),
			}
		}
	}
	return nameResolvers, err
}

//go:linkname parseNameServer github.com/Dreamacro/clash/config.parseNameServer
func parseNameServer(servers []string) ([]dns.NameServer, error)
