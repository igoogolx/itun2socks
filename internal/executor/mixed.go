package executor

import (
	"fmt"
)

type MixedProxyDetail struct {
	HubAddress string `json:"hubAddress"`
}

type MixedProxyClient struct {
	sysClient *SystemProxyClient
	tunClient *TunClient
}

func (m *MixedProxyClient) RuntimeDetail(hubAddress string) (interface{}, error) {
	return m.tunClient.RuntimeDetail(hubAddress)
}

func (m *MixedProxyClient) Start() error {
	sysErr := m.sysClient.Start()
	tunErr := m.tunClient.Start()

	if sysErr != nil && tunErr != nil {
		return fmt.Errorf("%v, %v", sysErr, tunErr)
	}
	if sysErr != nil {
		return sysErr
	}
	if tunErr != nil {
		return tunErr
	}

	return nil
}

func (m *MixedProxyClient) Close() error {
	sysErr := m.sysClient.Close()
	tunErr := m.tunClient.Close()

	if sysErr != nil && tunErr != nil {
		return fmt.Errorf("%v, %v", sysErr, tunErr)
	}
	if sysErr != nil {
		return sysErr
	}
	if tunErr != nil {
		return tunErr
	}

	return nil
}
