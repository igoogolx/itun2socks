package resolver

import (
	"github.com/igoogolx/itun2socks/pkg/clash/component/system_dns"
	_ "github.com/metacubex/mihomo/config" //init dns.ParseNameServer
	metaC "github.com/metacubex/mihomo/constant"
	"github.com/metacubex/mihomo/dns"
)

func New(mainServer []string, proxyAdapter metaC.ProxyAdapter, defaultInterfaceName string, disableCache bool) (dns.Resolvers, error) {
	nameservers, err := parse(mainServer, defaultInterfaceName)
	if err != nil {
		return dns.Resolvers{}, err
	}

	for index, nameserver := range nameservers {

		nameserver.ProxyAdapter = proxyAdapter
		nameservers[index] = nameserver
	}

	cacheMaxSize := 0
	if disableCache {
		//It's a little trick to disable cache because 0 means 4096 internally.
		cacheMaxSize = 1
	}

	mainDnsClient := dns.NewResolver(dns.Config{
		Main:         nameservers,
		CacheMaxSize: cacheMaxSize,
	})

	return mainDnsClient, nil
}

func parse(servers []string, defaultInterfaceName string) ([]dns.NameServer, error) {
	nameResolvers, err := dns.ParseNameServer(servers)
	if err != nil {
		return nil, err
	}
	needSystemDns := false
	for index, nameResolver := range nameResolvers {
		//FIXME: remove dhcp
		if nameResolver.Net == "system" || (nameResolver.Net == "dhcp" && nameResolver.Addr == "auto") {
			nameResolvers[index] = dns.NameServer{
				Net: "system",
			}
			needSystemDns = true
		}
	}

	if needSystemDns {

		nameResolvers = append(nameResolvers, dns.NameServer{Net: "dhcp", Addr: defaultInterfaceName})

		systemDnsServers, resolveSystemDnsErr := system_dns.ResolverV4Servers(defaultInterfaceName)
		if resolveSystemDnsErr == nil {

			for _, server := range systemDnsServers {
				nameResolvers = append(nameResolvers, dns.NameServer{
					Net:  "udp",
					Addr: server,
				})
			}

		}
	}

	return nameResolvers, err
}
