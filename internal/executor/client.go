package executor

import (
	"fmt"
	"github.com/Dreamacro/clash/component/iface"
	"github.com/igoogolx/itun2socks/internal/conn"
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
	BoostDns                string   `json:"boostDns"`
	ProxyServer             string   `json:"proxyServer"`
}

type Client struct {
	sync.RWMutex
	tun                     sTun.Tun
	stack                   sTun.Stack
	localserver             localserver.Server
	defaultInterfaceHandler network_iface.Handler
	deviceName              string
	localDns                []string
	remoteDns               []string
	boostDns                string
	runtimeDetail           Detail
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
		TunInterfaceName:        c.deviceName,
		LocalDns:                c.localDns,
		RemoteDns:               c.remoteDns,
		BoostDns:                c.boostDns,
		ProxyServer:             conn.GetProxyAddr(),
	}, nil
}

func (c *Client) Start() error {
	var err error
	if err = c.stack.Start(); err != nil {
		return fmt.Errorf("fail to start stack: %v", err)
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
	if err = c.localserver.Stop(); err != nil {
		return err
	}
	return nil
}
