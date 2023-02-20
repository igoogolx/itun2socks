package executor

import (
	"context"
	"github.com/igoogolx/itun2socks/cfg"
	network_iface "github.com/igoogolx/itun2socks/components/network-iface"
	"github.com/igoogolx/itun2socks/components/proxy-handler"
	"github.com/igoogolx/itun2socks/configuration"
	"github.com/igoogolx/itun2socks/conn"
	"github.com/igoogolx/itun2socks/dns"
	localserver "github.com/igoogolx/itun2socks/local-server"
	"github.com/igoogolx/itun2socks/runtime-detail"
	"github.com/igoogolx/itun2socks/tunnel"
	sTun "github.com/sagernet/sing-tun"
	"net/netip"
	"time"
)

func New() (*Client, error) {
	rawConfig, err := configuration.Read()
	if err != nil {
		return nil, err
	}
	config, err := cfg.New(rawConfig)
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
	tun, err := sTun.Open(tunOptions)
	if err != nil {
		return nil, err
	}
	stack, err := sTun.NewStack("gvisor", sTun.StackOptions{
		Context:    context.TODO(),
		Handler:    proxy_handler.New(),
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
	}, nil
}

func updateCfg(config cfg.Config) error {
	updateMatcher(config)
	updateConn(config)
	if err := updateHub(config); err != nil {
		return err
	}
	updateDns(config)
	return nil
}

func updateMatcher(c cfg.Config) {
	tunnel.UpdateMatcher(c.Rule)
}

func updateDns(c cfg.Config) {
	dns.UpdateMatcher(c.Rule)
}

func updateHub(c cfg.Config) error {
	return runtime_detail.Update(c.Device.Name, c.Rule.Dns.Local.Client.Nameservers(), c.Rule.Dns.Remote.Client.Nameservers())
}

func updateConn(c cfg.Config) {
	conn.UpdateProxy(c.Proxy)
}
