package executor

import (
	"context"
	"fmt"
	"net/netip"
	"time"

	cResolver "github.com/Dreamacro/clash/component/resolver"
	"github.com/igoogolx/itun2socks/internal/cfg"
	"github.com/igoogolx/itun2socks/internal/cfg/distribution/rule_engine"
	"github.com/igoogolx/itun2socks/internal/cfg/local_server"
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
)

type Client interface {
	Start() error
	Close() error
	RuntimeDetail(hubAddress string) (interface{}, error)
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
	rEngine, err := rule_engine.New(selectedRule, rawConfig.Rules)
	if err != nil {
		return "", err
	}
	matcher.UpdateRuleEngine(rEngine)
	log.Infoln(log.FormatLog(log.ExecutorPrefix, "update rule: %v"), selectedRule)
	dns.ResetCache()
	return selectedRule, nil
}

func newTun(isLocalServerEnabled bool) (*TunClient, error) {
	err := network_iface.StartMonitor()
	if err != nil {
		return nil, err
	}

	for {
		defaultInterface := network_iface.GetDefaultInterfaceName()
		if len(defaultInterface) != 0 {
			break
		}
		log.Infoln("%s", log.FormatLog(log.InitPrefix, "waiting for default interface name"))
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
	stack, err := sTun.NewStack("gvisor", sTun.StackOptions{
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

	newLocalServer := localserver.NewListener(config.LocalServer.Addr, config.LocalServer.Port)
	var matchers = []conn.Matcher{
		config.Rule.ConnMatcher,
	}
	if config.BlockQuic {
		matchers = append(matchers, conn.RejectQuicMather)
	}

	tunnel.UpdateShouldFindProcess(config.ShouldFindProcess)
	conn.UpdateConnMatcher(matchers)
	conn.UpdateIsFakeIpEnabled(config.FakeIp)
	conn.UpdateProxy(config.Proxy)

	log.Infoln(log.FormatLog(log.ExecutorPrefix, "set proxy: %v"), config.Proxy.Name())
	dns.UpdateDnsMap(config.Rule.Dns.Local.Client, config.Rule.Dns.Remote.Client)
	log.Infoln(log.FormatLog(log.ExecutorPrefix, "set dns, local: %v, remote: %v"), config.Rule.Dns.Local.Addresses, config.Rule.Dns.Remote.Addresses)
	_, err = UpdateRule()
	if err != nil {
		return nil, err
	}

	return &TunClient{
		stack:                stack,
		tun:                  tun,
		localserver:          newLocalServer,
		config:               config,
		isLocalServerEnabled: isLocalServerEnabled,
	}, nil
}

func newSysProxy() (*SystemProxyClient, error) {
	config, err := cfg.NewSystemProxy()
	if err != nil {
		return nil, err
	}

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

	newLocalServer := localserver.NewListener(config.LocalServer.Addr, config.LocalServer.Port)
	return &SystemProxyClient{
		localserver:     newLocalServer,
		activeInterface: config.ActiveInterface,
	}, nil
}

func newMixed() (Client, error) {
	rawConfig, err := configuration.Read()
	if err != nil {
		return nil, err
	}
	localServerConfig := local_server.New(rawConfig.Setting.LocalServer)
	newLocalServer := localserver.NewListener(localServerConfig.Addr, localServerConfig.Port)
	sysClient := &SystemProxyClient{
		localserver:     newLocalServer,
		activeInterface: rawConfig.Setting.HijackDns.NetworkService,
	}

	tunClient, err := newTun(false)
	if err != nil {
		return nil, err
	}
	return &MixedProxyClient{
		sysClient: sysClient,
		tunClient: tunClient,
	}, nil
}

func New() (Client, error) {
	cResolver.DefaultResolver = nil
	rawConfig, err := configuration.Read()
	if err != nil {
		return nil, err
	}
	if rawConfig.Setting.Mode == "tun" {
		return newTun(true)
	}
	if rawConfig.Setting.Mode == "system" {
		return newSysProxy()
	}
	if rawConfig.Setting.Mode == "mixed" {
		return newMixed()
	}
	return nil, fmt.Errorf("invalid proxy mode: %v", rawConfig.Setting.Mode)
}
