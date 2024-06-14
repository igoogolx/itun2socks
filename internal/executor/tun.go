package executor

import (
	"fmt"
	"github.com/Dreamacro/clash/component/iface"
	"github.com/igoogolx/itun2socks/internal/cfg"
	"github.com/igoogolx/itun2socks/internal/dns"
	localserver "github.com/igoogolx/itun2socks/internal/local_server"
	"github.com/igoogolx/itun2socks/internal/tunnel/statistic"
	"github.com/igoogolx/itun2socks/pkg/network_iface"
	sTun "github.com/sagernet/sing-tun"
	"sync"
)

type Detail struct {
	DirectedInterfaceName   string   `json:"directedInterfaceName"`
	DirectedInterfaceV4Addr string   `json:"directedInterfaceV4Addr"`
	TunInterfaceName        string   `json:"tunInterfaceName"`
	LocalDns                []string `json:"localDns"`
	RemoteDns               []string `json:"remoteDns"`
	BoostDns                []string `json:"boostDns"`
}

type TunClient struct {
	sync.RWMutex
	tun         sTun.Tun
	stack       sTun.Stack
	localserver localserver.Listener
	config      *cfg.Config
}

func (c *TunClient) RuntimeDetail() (interface{}, error) {
	networkInterface, err := iface.ResolveInterface(network_iface.GetDefaultInterfaceName())
	if err != nil {
		return nil, err
	}
	addr, err := networkInterface.PickIPv4Addr(nil)
	if err != nil {
		return nil, err
	}
	return &Detail{
		DirectedInterfaceV4Addr: addr.IP.String(),
		DirectedInterfaceName:   networkInterface.Name,
		TunInterfaceName:        c.config.Device.Name,
		LocalDns:                c.config.Rule.Dns.Local.Addresses,
		RemoteDns:               c.config.Rule.Dns.Remote.Addresses,
		BoostDns:                c.config.Rule.Dns.Boost.Addresses,
	}, nil
}

func (c *TunClient) Start() error {
	var err error
	if err = c.stack.Start(); err != nil {
		return fmt.Errorf("fail to start stack: %v", err)
	}
	if c.config.HijackDns.Enabled {

		err := dns.Hijack(c.config.HijackDns.NetworkService)
		if err != nil {
			return err
		}
	}
	err = c.localserver.Start()
	if err != nil {
		return err
	}
	return nil
}

func (c *TunClient) Close() error {
	var err error

	statistic.DefaultManager.CloseAllConnections()
	if err = c.tun.Close(); err != nil {
		return err
	}
	err = network_iface.StopMonitor()
	if err != nil {
		return err
	}
	if c.config.HijackDns.Enabled {
		err := dns.Resume(c.config.HijackDns.NetworkService)
		if err != nil {
			return err
		}
	}
	if err = c.localserver.Close(); err != nil {
		return err
	}
	return nil
}
