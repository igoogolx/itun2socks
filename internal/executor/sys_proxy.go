package executor

import (
	"github.com/igoogolx/itun2socks/internal/cfg"
	localserver "github.com/igoogolx/itun2socks/internal/local_server"
	"github.com/igoogolx/itun2socks/internal/tunnel/statistic"
	"sync"
)

type SysProxyDetail struct {
}

type SystemProxyClient struct {
	sync.RWMutex
	localserver localserver.Listener
	config      *cfg.Config
}

func (c *SystemProxyClient) RuntimeDetail() (interface{}, error) {
	return &Detail{}, nil
}

func (c *SystemProxyClient) Start() error {
	var err error
	err = c.localserver.Start()
	if err != nil {
		return err
	}
	return nil
}

func (c *SystemProxyClient) Close() error {
	var err error
	statistic.DefaultManager.CloseAllConnections()
	if err = c.localserver.Close(); err != nil {
		return err
	}
	return nil
}
