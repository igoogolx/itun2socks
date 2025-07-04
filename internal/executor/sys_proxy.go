package executor

import (
	"fmt"
	localserver "github.com/igoogolx/itun2socks/internal/local_server"
	"github.com/igoogolx/itun2socks/internal/tunnel/statistic"
	"github.com/igoogolx/itun2socks/pkg/sysproxy"
)

type SysProxyDetail struct {
	HubAddress string `json:"hubAddress"`
}

type SystemProxyClient struct {
	localserver localserver.Listener
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
	statistic.DefaultManager.CloseAllConnections()

	sysErr := sysproxy.Clear()

	lcErr := c.localserver.Close()

	if sysErr != nil && lcErr != nil {
		return fmt.Errorf("%v, %v", sysErr, lcErr)
	}
	if sysErr != nil {
		return sysErr
	}
	if lcErr != nil {
		return lcErr
	}

	return nil
}
