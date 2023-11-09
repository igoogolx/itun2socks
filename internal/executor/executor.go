package executor

import (
	"context"
	"github.com/igoogolx/itun2socks/internal/cfg"
	"github.com/igoogolx/itun2socks/internal/conn"
	"github.com/igoogolx/itun2socks/internal/dns"
	localserver "github.com/igoogolx/itun2socks/internal/local_server"
	"github.com/igoogolx/itun2socks/internal/proxy_handler"
	"github.com/igoogolx/itun2socks/internal/tunnel"
	"github.com/igoogolx/itun2socks/pkg/network_iface"
	sTun "github.com/sagernet/sing-tun"
	"net/netip"
	"time"
)

func New() (*Client, error) {
	err := network_iface.StartMonitor()
	if err != nil {
		return nil, err
	}
	config, err := cfg.New(network_iface.GetDefaultInterfaceName())
	if err != nil {
		return nil, err
	}
	tunOptions := sTun.Options{
		Name:         config.Device.Name,
		MTU:          uint32(config.Device.Mtu),
		Inet4Address: []netip.Prefix{netip.MustParsePrefix(config.Device.Gateway.String())},
		AutoRoute:    true,
		StrictRoute:  true,
	}
	tun, err := sTun.New(tunOptions)
	if err != nil {
		return nil, err
	}
	stack, err := sTun.NewStack("gvisor", sTun.StackOptions{
		Context:    context.TODO(),
		Handler:    proxy_handler.New(tunnel.TcpQueue(), tunnel.UdpQueue()),
		Tun:        tun,
		Name:       config.Device.Name,
		MTU:        uint32(config.Device.Mtu),
		UDPTimeout: int64(5 * time.Minute),
	})
	if err != nil {
		return nil, err
	}

	newLocalServer := localserver.NewListener(config.LocalServer.HttpAddr)

	updateCfg(*config)
	return &Client{
		stack:       stack,
		tun:         tun,
		localserver: newLocalServer,
		config:      config,
	}, nil
}

func updateCfg(config cfg.Config) {
	updateMatcher(config)
	updateConn(config)
	updateDns(config)
}

func updateMatcher(c cfg.Config) {
	conn.UpdateMatcher(c.Rule)
}

func updateDns(c cfg.Config) {
	dns.UpdateMatcher(c.Rule)
}

func updateConn(c cfg.Config) {
	conn.UpdateProxy(c.Proxy)
}
