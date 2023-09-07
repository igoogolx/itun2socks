package resolver

import (
	"context"
	cResolver "github.com/Dreamacro/clash/component/resolver"
	"github.com/Dreamacro/clash/dns"
	D "github.com/miekg/dns"
	"strings"
)

type Client struct {
	*dns.Resolver
	withAnswer func(msg *D.Msg)
}

func (c Client) ExchangeContext(ctx context.Context, m *D.Msg) (msg *D.Msg, err error) {
	if err == nil {
		c.withAnswer(m)
	}
	return c.ExchangeContext(ctx, m)
}

func New(mainServer, defaultServer string, withAnswer func(msg *D.Msg)) (cResolver.Resolver, error) {
	var boostNameResolver []dns.NameServer = nil
	if len(defaultServer) != 0 {
		boostNameResolver = []dns.NameServer{
			{
				Net:  "tcp",
				Addr: defaultServer,
			},
		}
	}
	localDnsNet := "tcp"
	if strings.Contains(mainServer, "https") {
		localDnsNet = "https"
	}
	localDnsClient := dns.NewResolver(dns.Config{
		Main: []dns.NameServer{{
			Net:  localDnsNet,
			Addr: mainServer,
		}},
		Default: boostNameResolver,
	})

	return Client{localDnsClient, withAnswer}, nil
}
