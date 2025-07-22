package executor

import (
	"fmt"
	localserver "github.com/igoogolx/itun2socks/internal/local_server"
	"github.com/igoogolx/itun2socks/internal/tunnel/statistic"
	"github.com/igoogolx/itun2socks/pkg/sysproxy"
	"strconv"
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

	addr := "127.0.0.1:" + strconv.Itoa(c.localserver.Port)
	err = sysproxy.Set(addr)
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
