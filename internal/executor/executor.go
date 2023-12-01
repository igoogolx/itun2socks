package executor

import (
	"context"
	"fmt"
	C "github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/cfg"
	"github.com/igoogolx/itun2socks/internal/configuration"
	"github.com/igoogolx/itun2socks/internal/conn"
	"github.com/igoogolx/itun2socks/internal/dns"
	localserver "github.com/igoogolx/itun2socks/internal/local_server"
	"github.com/igoogolx/itun2socks/internal/matcher"
	"github.com/igoogolx/itun2socks/internal/proxy_handler"
	"github.com/igoogolx/itun2socks/internal/tunnel"
	"github.com/igoogolx/itun2socks/pkg/network_iface"
	sTun "github.com/sagernet/sing-tun"
	"net/netip"
	"time"
)

type Client interface {
	Start() error
	Close() error
	RuntimeDetail() (interface{}, error)
}

func newTun() (Client, error) {
	err := network_iface.StartMonitor()
	if err != nil {
		return nil, err
	}
	config, err := cfg.NewTun(network_iface.GetDefaultInterfaceName())
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
		UDPTimeout: int64(5 * time.Second),
	})
	if err != nil {
		return nil, err
	}

	newLocalServer := localserver.NewListener(config.LocalServer.Addr)

	conn.UpdateConnMatcher([]conn.Matcher{
		config.Rule.ConnMatcher,
		conn.RejectQuicMather,
	})
	matcher.UpdateDnsMatcher(config.Rule)
	conn.UpdateProxy(config.Proxy)
	dns.UpdateDnsMap(config.Rule.Dns.Local.Client, config.Rule.Dns.Remote.Client, config.Rule.Dns.Boost.Client)

	return &TunClient{
		stack:       stack,
		tun:         tun,
		localserver: newLocalServer,
		config:      config,
	}, nil
}

func newSysProxy() (Client, error) {
	config, err := cfg.NewSystemProxy()
	if err != nil {
		return nil, err
	}
	newLocalServer := localserver.NewListener(config.LocalServer.Addr)
	conn.UpdateConnMatcher([]conn.Matcher{
		config.Rule.ConnMatcher,
	})
	matcher.UpdateDnsMatcher(config.Rule)
	conn.UpdateProxy(config.Proxy)
	return &SystemProxyClient{
		localserver: newLocalServer,
		config:      config,
	}, nil
}

func New() (Client, error) {
	rawConfig, err := configuration.Read()
	if err != nil {
		return nil, err
	}
	if rawConfig.Setting.Mode == "tun" {
		return newTun()
	}
	if rawConfig.Setting.Mode == "system" {
		return newSysProxy()
	}
	return nil, fmt.Errorf("invalid proxy mode: %v", rawConfig.Setting.Mode)
}

type RuntimeConfig struct {
	dnsMatcher matcher.Dns
	C.Proxy
}
