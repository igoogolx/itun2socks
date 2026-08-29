package resolver

import (
	_ "github.com/metacubex/mihomo/config" //init dns.ParseNameServer
	metaC "github.com/metacubex/mihomo/constant"
	"github.com/metacubex/mihomo/dns"
)

func New(mainServer []string, proxyAdapter metaC.ProxyAdapter, disableCache bool) (dns.Resolvers, error) {
	nameservers, err := parse(mainServer)
	if err != nil {
		return dns.Resolvers{}, err
	}

	for _, nameserver := range nameservers {

		nameserver.ProxyAdapter = proxyAdapter
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

func parse(servers []string) ([]dns.NameServer, error) {
	nameResolvers, err := dns.ParseNameServer(servers)
	if err != nil {
		return nil, err
	}
	for index, nameResolver := range nameResolvers {
		//FIXME: remove dhcp
		if nameResolver.Net == "system" || nameResolver.Net == "dhcp" {
			nameResolvers[index] = dns.NameServer{
				Net: "system",
			}
		}
	}
	return nameResolvers, err
}
