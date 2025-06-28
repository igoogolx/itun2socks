package executor

import (
	"sync"
)

type MixedProxyDetail struct {
	HubAddress string `json:"hubAddress"`
}

type MixedProxyClient struct {
	sync.RWMutex
	sysClient *SystemProxyClient
	tunClient *TunClient
}

func (m *MixedProxyClient) RuntimeDetail(hubAddress string) (interface{}, error) {
	return &MixedProxyDetail{hubAddress}, nil
}

func (m *MixedProxyClient) Start() error {
	var err error
	err = m.tunClient.Start()
	if err != nil {
		return err
	}
	err = m.sysClient.Start()
	if err != nil {
		return err
	}
	return nil
}

func (m *MixedProxyClient) Close() error {
	var err error
	err = m.tunClient.Close()
	if err != nil {
		return err
	}
	err = m.sysClient.Close()
	if err != nil {
		return err
	}
	return nil
}
