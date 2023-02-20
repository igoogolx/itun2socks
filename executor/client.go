package executor

import (
	"fmt"
	network_iface "github.com/igoogolx/itun2socks/components/network-iface"
	"github.com/igoogolx/itun2socks/dns"
	localserver "github.com/igoogolx/itun2socks/local-server"
	"github.com/igoogolx/itun2socks/tunnel/statistic"
	sTun "github.com/sagernet/sing-tun"
	"sync"
)

type Client struct {
	sync.RWMutex
	tun                     sTun.Tun
	stack                   sTun.Stack
	localserver             localserver.Server
	defaultInterfaceHandler network_iface.Handler
}

func (c *Client) Start() error {
	var err error
	if err = c.stack.Start(); err != nil {
		return fmt.Errorf("fail to start stack: %v", err)
	}
	if err = dns.FlushSysCaches(); err != nil {
		return fmt.Errorf("fail to flush dns cache: %v", err)
	}
	c.localserver.Start()
	return nil
}

func (c *Client) Close() error {
	var err error
	statistic.DefaultManager.CloseAllConnections()
	if err = c.stack.Close(); err != nil {
		return err
	}
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
