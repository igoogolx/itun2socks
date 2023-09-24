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
	interfaceHandler, err := network_iface.New()
	if err != nil {
		return nil, err
	}
	err = interfaceHandler.Monitor.Start()
	if err != nil {
		return nil, err
	}
	config, err := cfg.New()
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

	newLocalServer := localserver.New(config.LocalServer.HttpAddr)

	if err = updateCfg(config); err != nil {
		return nil, err
	}
	return &Client{
		stack:                   stack,
		tun:                     tun,
		localserver:             newLocalServer,
		defaultInterfaceHandler: *interfaceHandler,
		config:                  config,
	}, nil
}

func updateCfg(config cfg.Config) error {
	updateMatcher(config)
	updateConn(config)
	updateDns(config)
	return nil
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
