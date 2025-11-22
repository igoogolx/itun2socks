package executor

import (
	"fmt"
	"strconv"

	localserver "github.com/igoogolx/itun2socks/internal/local_server"
	"github.com/igoogolx/itun2socks/internal/tunnel/statistic"
	"github.com/igoogolx/itun2socks/pkg/network_iface"
	"github.com/igoogolx/itun2socks/pkg/sysproxy"
)

type SysProxyDetail struct {
	HubAddress              string `json:"hubAddress"`
	DirectedInterfaceV4Addr string `json:"directedInterfaceV4Addr"`
}

type SystemProxyClient struct {
	localserver     localserver.Listener
	activeInterface string
}

func (c *SystemProxyClient) RuntimeDetail(hubAddress string) (interface{}, error) {
	return &SysProxyDetail{hubAddress, network_iface.GetLanV4Address()}, nil
}

func (c *SystemProxyClient) Start() error {
	var err error
	err = c.localserver.Start()
	if err != nil {
		return err
	}

	addr := "127.0.0.1:" + strconv.Itoa(c.localserver.Port)
	err = sysproxy.Set(addr, c.activeInterface)
	if err != nil {
		return err
	}
	return nil
}

func (c *SystemProxyClient) Close() error {
	statistic.DefaultManager.CloseAllConnections()

	sysErr := sysproxy.Clear(c.activeInterface)

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
