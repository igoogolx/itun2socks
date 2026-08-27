package dns

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/metacubex/mihomo/component/dialer"
	"math/rand"
	"net"
	"net/netip"
	"strconv"
	"strings"

	metaC "github.com/metacubex/mihomo/constant"

	"github.com/igoogolx/itun2socks/pkg/clash/component/resolver"

	D "github.com/miekg/dns"
)

type client struct {
	*D.Client
	r         *Resolver
	port      string
	host      string
	iface     string
	getDialer func() (metaC.Proxy, error)
}

func (c *client) GetServers() []string {
	return []string{c.host}
}

func (c *client) Exchange(m *D.Msg) (*D.Msg, error) {
	return c.ExchangeContext(context.Background(), m)
}

func (c *client) ExchangeContext(ctx context.Context, m *D.Msg) (*D.Msg, error) {
	var (
		ip  netip.Addr
		err error
	)
	if c.r == nil {
		// a default ip dns
		if ip, err = netip.ParseAddr(c.host); err != nil {
			return nil, fmt.Errorf("dns %s not a valid ip", c.host)
		}
	} else {
		ips, err := resolver.LookupIPWithResolver(ctx, c.host, c.r)
		if err != nil {
			return nil, fmt.Errorf("use default dns resolve failed: %w", err)
		} else if len(ips) == 0 {
			return nil, fmt.Errorf("%w: %s", resolver.ErrIPNotFound, c.host)
		}

		ip = ips[rand.Intn(len(ips))]

	}

	network := metaC.UDP
	if strings.HasPrefix(c.Client.Net, "tcp") {
		network = metaC.TCP
	}

	numPort, err := strconv.ParseUint(c.port, 10, 16)

	if err != nil {
		return nil, err
	}
	var conn net.Conn
	connDial, err := c.getDialer()
	if err != nil {
		return nil, err
	}
	if network == metaC.TCP {
		conn, err = connDial.DialContext(ctx, &metaC.Metadata{
			NetWork: network,
			SrcIP:   netip.Addr{},
			DstIP:   ip,
			SrcPort: 0,
			DstPort: uint16(numPort),
			Host:    "",
		})
	} else {
		conn, err = dialer.DialContext(ctx, "udp", net.JoinHostPort(ip.String(), c.port))
	}

	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// miekg/dns ExchangeContext doesn't respond to context cancel.
	// this is a workaround
	type result struct {
		msg *D.Msg
		err error
	}
	ch := make(chan result, 1)
	go func() {
		if strings.HasSuffix(c.Client.Net, "tls") {
			conn = tls.Client(conn, c.Client.TLSConfig)
		}

		msg, _, err := c.Client.ExchangeWithConn(m, &D.Conn{
			Conn:         conn,
			UDPSize:      c.Client.UDPSize,
			TsigSecret:   c.Client.TsigSecret,
			TsigProvider: c.Client.TsigProvider,
		})

		ch <- result{msg, err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case ret := <-ch:
		return ret.msg, ret.err
	}
}
