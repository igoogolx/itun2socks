package network_iface

import (
	"context"
	"github.com/Dreamacro/clash/log"
	"github.com/igoogolx/itun2socks/configuration"
	tun "github.com/sagernet/sing-tun"
	E "github.com/sagernet/sing/common/exceptions"
	"go.uber.org/atomic"
	"net/netip"
)

var defaultInterfaceName = atomic.NewString("")

func GetDefaultInterfaceName() string {
	return defaultInterfaceName.Load()
}

type ErrorHandler struct {
}

func (e ErrorHandler) NewError(ctx context.Context, err error) {
	log.Errorln("network interface monitor: %v", err)
}

type Handler struct {
	Monitor tun.DefaultInterfaceMonitor
}

func New() (*Handler, error) {
	networkUpdateMonitor, err := tun.NewNetworkUpdateMonitor(ErrorHandler{})
	if err != nil {
		err = E.Cause(err, "create NetworkUpdateMonitor")
		return nil, err
	}
	err = networkUpdateMonitor.Start()
	if err != nil {
		err = E.Cause(err, "start NetworkUpdateMonitor")
		return nil, err
	}

	defaultInterfaceMonitor, err := tun.NewDefaultInterfaceMonitor(networkUpdateMonitor, tun.DefaultInterfaceMonitorOptions{OverrideAndroidVPN: true})
	if err != nil {
		err = E.Cause(err, "create DefaultInterfaceMonitor")
		return nil, err
	}
	defaultInterfaceMonitor.RegisterCallback(func(event int) error {
		return Update(defaultInterfaceMonitor.DefaultInterfaceName(netip.Addr{}))
	})
	return &Handler{defaultInterfaceMonitor}, nil
}

func Update(name string) error {
	setting, err := configuration.GetSetting()
	if err != nil {
		return err
	}
	nextName := name
	if len(setting.DefaultInterface) != 0 {
		nextName = setting.DefaultInterface
	}
	defaultInterfaceName.Store(nextName)
	return err
}
