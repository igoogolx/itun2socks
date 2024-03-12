package network_iface

import (
	"context"
	"github.com/Dreamacro/clash/component/dialer"
	"github.com/igoogolx/itun2socks/internal/configuration"
	"github.com/igoogolx/itun2socks/pkg/log"
	tun "github.com/sagernet/sing-tun"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sirupsen/logrus"
	"go.uber.org/atomic"
	"net/netip"
)

var defaultInterfaceName = atomic.NewString("")
var defaultInterfaceMonitor tun.DefaultInterfaceMonitor

func GetDefaultInterfaceName() string {
	return defaultInterfaceName.Load()
}

type ErrorHandler struct {
}

func (e ErrorHandler) NewError(_ context.Context, err error) {
	log.Errorln(log.FormatLog(log.ExecutorPrefix, "network interface monitor: %v"), err)
}

type Handler struct {
	Monitor tun.DefaultInterfaceMonitor
}

func StartMonitor() error {
	setting, err := configuration.GetSetting()
	if err != nil {
		return err
	}
	if len(setting.DefaultInterface) != 0 {
		update(setting.DefaultInterface)
		return nil
	}
	networkUpdateMonitor, err := tun.NewNetworkUpdateMonitor(logrus.StandardLogger())
	if err != nil {
		err = E.Cause(err, "create NetworkUpdateMonitor")
		return err
	}
	err = networkUpdateMonitor.Start()
	if err != nil {
		err = E.Cause(err, "start NetworkUpdateMonitor")
		return err
	}

	defaultInterfaceMonitor, err = tun.NewDefaultInterfaceMonitor(networkUpdateMonitor, logrus.StandardLogger(), tun.DefaultInterfaceMonitorOptions{OverrideAndroidVPN: true})
	if err != nil {
		err = E.Cause(err, "create DefaultInterfaceMonitor")
		return err
	}
	defaultInterfaceMonitor.RegisterCallback(func(event int) {
		update(defaultInterfaceMonitor.DefaultInterfaceName(netip.IPv4Unspecified()))
	})
	err = defaultInterfaceMonitor.Start()
	if err != nil {
		return err
	}
	update(defaultInterfaceMonitor.DefaultInterfaceName(netip.IPv4Unspecified()))
	return nil
}

func StopMonitor() error {
	defer func() {
		defaultInterfaceMonitor = nil
	}()
	if defaultInterfaceMonitor != nil {
		return defaultInterfaceMonitor.Close()
	}
	return nil
}

func update(name string) {
	defaultInterfaceName.Store(name)
	dialer.DefaultInterface.Store(name)
	log.Infoln(log.FormatLog(log.ExecutorPrefix, "update default interface: %v"), name)
}
