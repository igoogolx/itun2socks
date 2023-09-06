package executor

import (
	"context"
	"github.com/igoogolx/itun2socks/internal/cfg"
	conn2 "github.com/igoogolx/itun2socks/internal/conn"
	"github.com/igoogolx/itun2socks/internal/dns"
	localserver "github.com/igoogolx/itun2socks/internal/local_server"
	"github.com/igoogolx/itun2socks/internal/proxy_handler"
	tunnel2 "github.com/igoogolx/itun2socks/internal/tunnel"
	"github.com/igoogolx/itun2socks/pkg/network_iface"
	sTun "github.com/sagernet/sing-tun"
	"net/netip"
	"time"
)

func New() (*Client, error) {
	config, err := cfg.New()
	if err != nil {
		return nil, err
	}
	tunDevice := config.Device
	tunOptions := sTun.Options{
		Name:         tunDevice.Name,
		MTU:          uint32(tunDevice.Mtu),
		Inet4Address: []netip.Prefix{netip.MustParsePrefix(tunDevice.Gateway.String())},
		AutoRoute:    true,
		StrictRoute:  true,
	}
	tun, err := sTun.New(tunOptions)
	if err != nil {
		return nil, err
	}
	stack, err := sTun.NewStack("gvisor", sTun.StackOptions{
		Context:    context.TODO(),
		Handler:    proxy_handler.New(tunnel2.TcpQueue(), tunnel2.UdpQueue()),
		Tun:        tun,
		Name:       tunDevice.Name,
		MTU:        uint32(tunDevice.Mtu),
		UDPTimeout: int64(5 * time.Minute),
	})
	if err != nil {
		return nil, err
	}

	newLocalServer := localserver.New(config.LocalServer.HttpAddr)
	interfaceHandler, err := network_iface.New()
	if err != nil {
		return nil, err
	}
	err = interfaceHandler.Monitor.Start()
	if err != nil {
		return nil, err
	}
	if err = updateCfg(config); err != nil {
		return nil, err
	}
	return &Client{
		stack:                   stack,
		tun:                     tun,
		localserver:             newLocalServer,
		defaultInterfaceHandler: *interfaceHandler,
		deviceName:              tunDevice.Name,
		localDns:                []string{config.Rule.Dns.Local.Address},
		remoteDns:               []string{config.Rule.Dns.Remote.Address},
		boostDns:                config.Rule.Dns.BoostNameserver,
	}, nil
}

func updateCfg(config cfg.Config) error {
	updateMatcher(config)
	updateConn(config)
	updateDns(config)
	return nil
}

func updateMatcher(c cfg.Config) {
	conn2.UpdateMatcher(c.Rule)
}

func updateDns(c cfg.Config) {
	dns.UpdateMatcher(c.Rule)
}

func updateConn(c cfg.Config) {
	conn2.UpdateProxy(c.Proxy)
}
