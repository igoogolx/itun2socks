package dns

import (
	metaResolver "github.com/metacubex/mihomo/component/resolver"
	"github.com/metacubex/mihomo/dns"
)

func InitDnsForSysProxy() error {

	nameServers, err := dns.ParseNameServer([]string{"system"})
	if err != nil {
		return err
	}

	metaResolver.SystemResolver = dns.NewResolver(dns.Config{Main: nameServers, CacheMaxSize: 1})

	return nil

}

func InitDnsForTunProxy(sysResolvers dns.Resolvers) {
	ResetCache()
	metaResolver.DefaultResolver = sysResolvers
	metaResolver.SystemResolver = sysResolvers
}
