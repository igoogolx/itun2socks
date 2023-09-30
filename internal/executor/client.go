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
	"runtime"
	"sync"
)

type Detail struct {
	DirectedInterfaceName   string   `json:"directedInterfaceName"`
	DirectedInterfaceV4Addr string   `json:"directedInterfaceV4Addr"`
	TunInterfaceName        string   `json:"tunInterfaceName"`
	LocalDns                []string `json:"localDns"`
	RemoteDns               []string `json:"remoteDns"`
	BoostDns                string   `json:"boostDns"`
}

type Client struct {
	sync.RWMutex
	tun                     sTun.Tun
	stack                   sTun.Stack
	localserver             localserver.Server
	defaultInterfaceHandler network_iface.Handler
	config                  cfg.Config
}

func (c *Client) RuntimeDetail() (*Detail, error) {
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
		LocalDns:                []string{c.config.Rule.Dns.Local.Address},
		RemoteDns:               []string{c.config.Rule.Dns.Remote.Address},
		BoostDns:                c.config.Rule.Dns.Boost.Address,
	}, nil
}

func (c *Client) Start() error {
	var err error
	if err = c.stack.Start(); err != nil {
		return fmt.Errorf("fail to start stack: %v", err)
	}
	if c.config.HijackDns.Enabled {
		if runtime.GOOS == "darwin" {
			err := dns.Hijack(c.config.HijackDns.NetworkService)
			if err != nil {
				return err
			}
		}
	}
	c.localserver.Start()
	return nil
}

func (c *Client) Close() error {
	var err error
	statistic.DefaultManager.CloseAllConnections()
	if err = c.tun.Close(); err != nil {
		return err
	}
	err = c.defaultInterfaceHandler.Monitor.Close()
	if err != nil {
		return err
	}
	if c.config.HijackDns.Enabled {
		if runtime.GOOS == "darwin" {
			err := dns.Resume(c.config.HijackDns.NetworkService)
			if err != nil {
				return err
			}
		}
	}
	if err = c.localserver.Stop(); err != nil {
		return err
	}
	return nil
}
