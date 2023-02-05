package cfg

import (
	"github.com/Dreamacro/clash/constant"
	lru "github.com/hashicorp/golang-lru"
	"github.com/igoogolx/itun2socks/cfg/distribution"
	"github.com/igoogolx/itun2socks/cfg/local-server"
	"github.com/igoogolx/itun2socks/cfg/outbound"
	"github.com/igoogolx/itun2socks/cfg/tun"
	db2 "github.com/igoogolx/itun2socks/configuration"
	"github.com/igoogolx/itun2socks/configuration/configuration-types"
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

func New(rawConfig configuration_types.Config) (Config, error) {
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
	proxy, err := outbound.New(rawConfig.Proxy, rawConfig.Selected.Proxy)
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
