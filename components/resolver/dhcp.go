package resolver

import (
	"context"
	"github.com/Dreamacro/clash/component/dhcp"
	"github.com/igoogolx/itun2socks/global"
	"github.com/miekg/dns"
	"net"
	"time"
)

var (
	DhcpResolveTimeout = 3 * time.Second
)

type dhcpClient struct {
}

func (d dhcpClient) Exchange(m *dns.Msg) (*dns.Msg, error) {
	return d.ExchangeContext(context.Background(), m)
}

func (d dhcpClient) ExchangeContext(ctx context.Context, m *dns.Msg) (*dns.Msg, error) {
	nameservers, err := d.resolveDns()
	if err != nil {
		return nil, err
	}
	client := NewResolver(nameservers)
	return client.ExchangeContext(ctx, m)
}

func (d dhcpClient) resolveDns() ([]net.IP, error) {
	rCtx, cancel := context.WithTimeout(context.Background(), DhcpResolveTimeout)
	defer cancel()
	return dhcp.ResolveDNSFromDHCP(rCtx, global.GetDefaultInterfaceName())
}

func (d dhcpClient) Nameservers() []string {
	nameservers, err := d.resolveDns()
	if err != nil {
		return []string{}
	}
	servers := make([]string, 0, len(nameservers))
	for _, n := range nameservers {
		servers = append(servers, n.String())
	}
	return servers
}

func newDhcpClient() dhcpClient {
	return dhcpClient{}
}
