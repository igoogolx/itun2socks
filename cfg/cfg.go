package cfg

import (
	"github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/cfg/distribution"
	"github.com/igoogolx/itun2socks/cfg/local-server"
	"github.com/igoogolx/itun2socks/cfg/outbound"
	"github.com/igoogolx/itun2socks/cfg/tun"
	db "github.com/igoogolx/itun2socks/configuration"
)

type Config struct {
	Rule        distribution.Config
	Proxy       constant.Proxy
	Device      tun.Config
	LocalServer local_server.Cfg
}

func New(rawConfig db.Config) (Config, error) {
	selectedRule, err := db.GetSelectedRule()
	if err != nil {
		return Config{}, err
	}
	rule, err := distribution.New(rawConfig.Setting.Dns.Boost, rawConfig.Setting.Dns.Remote, rawConfig.Setting.Dns.Local, selectedRule, rawConfig.Setting.TrueProxyServer)
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
