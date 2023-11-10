package cfg

import (
	"github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/cfg/distribution"
	"github.com/igoogolx/itun2socks/internal/cfg/local_server"
	"github.com/igoogolx/itun2socks/internal/cfg/outbound"
	"github.com/igoogolx/itun2socks/internal/configuration"
)

type SystemProxyConfig struct {
	Rule        distribution.SystemProxyConfig
	Proxy       constant.Proxy
	LocalServer local_server.Cfg
}

func NewSystemProxy() (*SystemProxyConfig, error) {
	rawConfig, err := configuration.Read()
	if err != nil {
		return nil, err
	}
	selectedRule, err := configuration.GetSelectedRule()
	if err != nil {
		return nil, err
	}

	rule, err := distribution.NewSystemProxy(
		selectedRule,
		rawConfig.Rules,
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
	return &SystemProxyConfig{
		rule,
		proxy,
		newLocalServer,
	}, nil
}
