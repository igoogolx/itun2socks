package cfg

import (
	"github.com/Dreamacro/clash/constant"
	lru "github.com/hashicorp/golang-lru"
	"github.com/igoogolx/itun2socks/cfg/distribution"
	"github.com/igoogolx/itun2socks/cfg/local-server"
	"github.com/igoogolx/itun2socks/cfg/outbound"
	"github.com/igoogolx/itun2socks/cfg/tun"
	db2 "github.com/igoogolx/itun2socks/configuration"
)

var (
	DnsTable, _ = lru.New(1000)
)

type Config struct {
	Rule        distribution.Config
	Proxy       constant.Proxy
	Device      tun.Config
	LocalServer local_server.Cfg
}

func New(rawConfig db2.Config) (Config, error) {
	selectedRule, err := db2.GetSelectedRule()
	if err != nil {
		return Config{}, err
	}
	rule, err := distribution.New(selectedRule, rawConfig.Setting.TrueProxyServer, DnsTable)
	if err != nil {
		return Config{}, err
	}
	device, err := tun.New()
	if err != nil {
		return Config{}, err
	}
	outboundOption := outbound.Option{
		Mode:    rawConfig.Setting.Outbound.Mode,
		Proxies: rawConfig.Proxy,
	}
	if rawConfig.Setting.Outbound.Mode == "select" {
		outboundOption.Config = map[string]string{
			"selected": rawConfig.Selected.Proxy,
		}
	} else if rawConfig.Setting.Outbound.Mode == "auto" {
		outboundOption.Config = rawConfig.Setting.Outbound.Config
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
