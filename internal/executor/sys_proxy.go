package executor

import (
	"github.com/getlantern/sysproxy"
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
	config      *cfg.SystemProxyConfig
	off         func() error
}

func (c *SystemProxyClient) RuntimeDetail() (interface{}, error) {
	return &SysProxyDetail{}, nil
}

func (c *SystemProxyClient) Start() error {
	var err error
	err = c.localserver.Start()
	if err != nil {
		return err
	}
	helperFullPath := "sysproxy-cmd"
	err = sysproxy.EnsureHelperToolPresent(helperFullPath, "Input your password and save the world!", "")
	if err != nil {
		return err
	}
	off, err := sysproxy.On(c.localserver.Addr)
	if err != nil {
		return err
	}
	c.off = off
	return nil
}

func (c *SystemProxyClient) Close() error {
	var err error
	statistic.DefaultManager.CloseAllConnections()
	err = c.off()
	if err != nil {
		return err
	}
	if err = c.localserver.Close(); err != nil {
		return err
	}

	return nil
}
