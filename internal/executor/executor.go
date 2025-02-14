package executor

import (
	"context"
	"fmt"
	cResolver "github.com/Dreamacro/clash/component/resolver"
	"github.com/igoogolx/itun2socks/internal/cfg"
	"github.com/igoogolx/itun2socks/internal/cfg/distribution/ruleEngine"
	"github.com/igoogolx/itun2socks/internal/configuration"
	"github.com/igoogolx/itun2socks/internal/conn"
	"github.com/igoogolx/itun2socks/internal/dns"
	localserver "github.com/igoogolx/itun2socks/internal/local_server"
	"github.com/igoogolx/itun2socks/internal/matcher"
	"github.com/igoogolx/itun2socks/internal/proxy_handler"
	"github.com/igoogolx/itun2socks/internal/tunnel"
	"github.com/igoogolx/itun2socks/pkg/log"
	"github.com/igoogolx/itun2socks/pkg/network_iface"
	sTun "github.com/sagernet/sing-tun"
	"github.com/sirupsen/logrus"
	"net/netip"
	"time"
)

type Client interface {
	Start() error
	Close() error
	RuntimeDetail() (interface{}, error)
}

func UpdateRule() (string, error) {
	rawConfig, err := configuration.Read()
	if err != nil {
		return "", err
	}
	selectedRule, err := configuration.GetSelectedRule()
	if err != nil {
		return "", err
	}
	rEngine, err := ruleEngine.New(selectedRule, rawConfig.Rules)
	if err != nil {
		return "", err
	}
	matcher.UpdateRule(rEngine)
	log.Infoln(log.FormatLog(log.ExecutorPrefix, "update rule: %v"), selectedRule)
	return selectedRule, nil
}

func newTun() (Client, error) {
	err := network_iface.StartMonitor()
	if err != nil {
		return nil, err
	}

	for {
		defaultInterface := network_iface.GetDefaultInterfaceName()
		if len(defaultInterface) != 0 {
			break
		}
		log.Infoln(log.FormatLog(log.InitPrefix, "waiting for default interface name"))
		time.Sleep(1 * time.Second)
	}

	config, err := cfg.NewTun(network_iface.GetDefaultInterfaceName())
	if err != nil {
		return nil, err
	}
	tunOptions := sTun.Options{
		Name:         config.Device.Name,
		MTU:          uint32(config.Device.Mtu),
		Inet4Address: []netip.Prefix{config.Device.Gateway},
		AutoRoute:    true,
		StrictRoute:  true,
		Logger:       logrus.StandardLogger(),
	}
	tun, err := sTun.New(tunOptions)
	if err != nil {
		return nil, err
	}
	err = tun.Start()
	if err != nil {
		return nil, err
	}
	stack, err := sTun.NewStack(config.Stack, sTun.StackOptions{
		Context:    context.Background(),
		Handler:    proxy_handler.New(tunnel.TcpQueue(), tunnel.UdpQueue()),
		TunOptions: tunOptions,
		Tun:        tun,
		UDPTimeout: 5 * time.Second,
		Logger:     logrus.StandardLogger(),
	})
	if err != nil {
		return nil, err
	}

	log.Infoln(log.FormatLog(log.ExecutorPrefix, "set network stack: %v"), config.Stack)

	newLocalServer := localserver.NewListener(config.LocalServer.Addr)
	var matchers = []conn.Matcher{
		config.Rule.ConnMatcher,
	}
	if config.BlockQuic {
		matchers = append(matchers, conn.RejectQuicMather)
	}
	tunnel.UpdateShouldFindProcess(config.ShouldFindProcess)
	conn.UpdateConnMatcher(matchers)

	conn.UpdateProxy(config.Proxy)
	log.Infoln(log.FormatLog(log.ExecutorPrefix, "set proxy: %v"), config.Proxy.Name())
	dns.UpdateDnsMap(config.Rule.Dns.Local.Client, config.Rule.Dns.Remote.Client)
	log.Infoln(log.FormatLog(log.ExecutorPrefix, "set dns, local: %v, remote: %v"), config.Rule.Dns.Local.Addresses, config.Rule.Dns.Remote.Addresses)
	_, err = UpdateRule()
	if err != nil {
		return nil, err
	}

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

	tunnel.UpdateShouldFindProcess(false)
	conn.UpdateConnMatcher([]conn.Matcher{
		config.Rule.ConnMatcher,
	})
	conn.UpdateProxy(config.Proxy)
	log.Infoln(log.FormatLog(log.ExecutorPrefix, "set proxy: %v"), config.Proxy.Name())
	_, err = UpdateRule()
	if err != nil {
		return nil, err
	}

	return &SystemProxyClient{
		localserver: newLocalServer,
		config:      config,
	}, nil
}

func New() (Client, error) {
	cResolver.DefaultResolver = nil
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
