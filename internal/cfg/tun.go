package cfg

import (
	"github.com/igoogolx/itun2socks/internal/cfg/distribution"
	"github.com/igoogolx/itun2socks/internal/cfg/local_server"
	"github.com/igoogolx/itun2socks/internal/cfg/outbound"
	"github.com/igoogolx/itun2socks/internal/cfg/tun"
	"github.com/igoogolx/itun2socks/internal/configuration"
	"github.com/igoogolx/itun2socks/internal/dns"
	"github.com/metacubex/mihomo/component/fakeip"
	metaC "github.com/metacubex/mihomo/constant"
)

type Config struct {
	Rule              distribution.Config
	Proxy             metaC.Proxy
	Device            *tun.Config
	LocalServer       local_server.Cfg
	HijackDns         configuration.HijackDns
	BlockQuic         bool
	Stack             string
	ShouldFindProcess bool
	FakeIp            bool
}

func NewTun(defaultInterfaceName string) (*Config, error) {
	rawConfig, err := configuration.Read()
	if err != nil {
		return nil, err
	}
	device, err := tun.New()
	if err != nil {
		return nil, err
	}
	disableDnsCache := rawConfig.Setting.Dns.DisableCache
	remoteDnsItems := rawConfig.Setting.Dns.Server.Remote
	var fakeIpPool *fakeip.Pool
	fakeIp := rawConfig.Setting.Dns.FakeIp
	if fakeIp {
		fakeIpPool = dns.FakeIpPool
	}
	rule, err := distribution.NewTun(
		rawConfig.Setting.Dns.Server.Boost,
		remoteDnsItems,
		rawConfig.Setting.Dns.Server.Local,
		disableDnsCache,
		fakeIpPool,
		defaultInterfaceName,
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
		fakeIp,
	}, nil
}
