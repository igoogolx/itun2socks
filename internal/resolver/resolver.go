package resolver

import (
	cResolver "github.com/Dreamacro/clash/component/resolver"
	"github.com/Dreamacro/clash/dns"
	"strings"
)

func New(mainServer, defaultServer string) (cResolver.Resolver, error) {
	var defaultNameResolver []dns.NameServer = nil
	if len(defaultServer) != 0 {
		defaultNameResolver = []dns.NameServer{
			{
				Net:  "tcp",
				Addr: defaultServer,
			},
		}
	}
	mainDnsNet := "tcp"
	if strings.Contains(mainServer, "https") {
		mainDnsNet = "https"
	}
	mainDnsClient := dns.NewResolver(dns.Config{
		Main: []dns.NameServer{{
			Net:  mainDnsNet,
			Addr: mainServer,
		}},
		Default: defaultNameResolver,
	})

	return mainDnsClient, nil
}
