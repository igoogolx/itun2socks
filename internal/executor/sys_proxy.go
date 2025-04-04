package executor

import (
	"github.com/igoogolx/itun2socks/internal/cfg"
	localserver "github.com/igoogolx/itun2socks/internal/local_server"
	"github.com/igoogolx/itun2socks/internal/tunnel/statistic"
	"github.com/igoogolx/itun2socks/pkg/sysproxy"
	"sync"
)

type SysProxyDetail struct {
	HubAddress string `json:"hubAddress"`
}

type SystemProxyClient struct {
	sync.RWMutex
	localserver localserver.Listener
	config      *cfg.SystemProxyConfig
}

func (c *SystemProxyClient) RuntimeDetail(hubAddress string) (interface{}, error) {
	return &SysProxyDetail{hubAddress}, nil
}

func (c *SystemProxyClient) Start() error {
	var err error
	err = c.localserver.Start()
	if err != nil {
		return err
	}
	err = sysproxy.Set(c.localserver.Addr)
	if err != nil {
		return err
	}
	return nil
}

func (c *SystemProxyClient) Close() error {
	var err error
	statistic.DefaultManager.CloseAllConnections()
	err = sysproxy.Clear()
	if err != nil {
		return err
	}
	if err = c.localserver.Close(); err != nil {
		return err
	}

	return nil
}
