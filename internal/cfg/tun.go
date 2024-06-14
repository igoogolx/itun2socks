package cfg

import (
	"github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/cfg/distribution"
	"github.com/igoogolx/itun2socks/internal/cfg/local_server"
	"github.com/igoogolx/itun2socks/internal/cfg/outbound"
	"github.com/igoogolx/itun2socks/internal/cfg/tun"
	"github.com/igoogolx/itun2socks/internal/configuration"
)

type Config struct {
	Rule              distribution.Config
	Proxy             constant.Proxy
	Device            *tun.Config
	LocalServer       local_server.Cfg
	HijackDns         configuration.HijackDns
	BlockQuic         bool
	Stack             string
	ShouldFindProcess bool
}

func NewTun(defaultInterfaceName string) (*Config, error) {
	rawConfig, err := configuration.Read()
	if err != nil {
		return nil, err
	}
	selectedRule, err := configuration.GetSelectedRule()
	if err != nil {
		return nil, err
	}
	device, err := tun.New()
	if err != nil {
		return nil, err
	}
	rule, err := distribution.NewTun(
		rawConfig.Setting.Dns.Server.Boost,
		rawConfig.Setting.Dns.Server.Remote,
		rawConfig.Setting.Dns.Server.Local,
		selectedRule,
		rawConfig.Rules,
		defaultInterfaceName,
		rawConfig.Setting.Dns.DisableCache,
	)

	if err != nil {
		return nil, err
	}
	outboundOption := outbound.Option{
		AutoMode:      rawConfig.Setting.AutoMode,
		Proxies:       rawConfig.Proxy,
		SelectedProxy: rawConfig.Selected.Proxy,
	}
	proxy, err := outbound.New(outboundOption)
	if err != nil {
		return nil, err
	}
	newLocalServer := local_server.New(rawConfig.Setting.LocalServer)

	return &Config{
		rule,
		proxy,
		device,
		newLocalServer,
		rawConfig.Setting.HijackDns,
		rawConfig.Setting.BlockQuic,
		rawConfig.Setting.Stack,
		rawConfig.Setting.ShouldFindProcess,
	}, nil
}
