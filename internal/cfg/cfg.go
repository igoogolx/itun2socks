package cfg

import (
	"github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/cfg/distribution"
	"github.com/igoogolx/itun2socks/internal/cfg/local-server"
	"github.com/igoogolx/itun2socks/internal/cfg/outbound"
	"github.com/igoogolx/itun2socks/internal/cfg/tun"
	"github.com/igoogolx/itun2socks/internal/configuration"
)

type Config struct {
	Rule        distribution.Config
	Proxy       constant.Proxy
	Device      tun.Config
	LocalServer local_server.Cfg
}

func New() (Config, error) {
	rawConfig, err := configuration.Read()
	if err != nil {
		return Config{}, err
	}
	selectedRule, err := configuration.GetSelectedRule()
	if err != nil {
		return Config{}, err
	}
	rule, err := distribution.New(rawConfig.Setting.Dns.Boost.Value, rawConfig.Setting.Dns.Remote.Value, rawConfig.Setting.Dns.Local.Value, selectedRule, rawConfig.Setting.TrueProxyServer)
	if err != nil {
		return Config{}, err
	}
	device, err := tun.New()
	if err != nil {
		return Config{}, err
	}
	outboundOption := outbound.Option{
		AutoMode:      rawConfig.Setting.AutoMode,
		Proxies:       rawConfig.Proxy,
		SelectedProxy: rawConfig.Selected.Proxy,
	}
	proxy, err := outbound.New(outboundOption)
	if err != nil {
		return Config{}, err
	}
	newLocalServer := local_server.New(rawConfig.Setting.LocalServer)
	return Config{
		rule,
		proxy,
		device,
		newLocalServer,
	}, nil
}