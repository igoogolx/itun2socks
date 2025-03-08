//go:build windows && arm64

package executor

import (
	"fmt"
	"github.com/igoogolx/itun2socks/internal/cfg"
	localserver "github.com/igoogolx/itun2socks/internal/local_server"
	"sync"
)

type SysProxyDetail struct {
	HubAddress string `json:"hubAddress"`
}

type SystemProxyClient struct {
	sync.RWMutex
	localserver localserver.Listener
	config      *cfg.SystemProxyConfig
	off         func() error
}

func (c *SystemProxyClient) RuntimeDetail(hubAddress string) (interface{}, error) {
	return &SysProxyDetail{hubAddress}, fmt.Errorf("not implemented yet")
}

func (c *SystemProxyClient) Start() error {
	return fmt.Errorf("not implemented yet")
}

func (c *SystemProxyClient) Close() error {
	return fmt.Errorf("not implemented yet")

}
