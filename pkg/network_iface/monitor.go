package network_iface

import (
	"context"
	"github.com/Dreamacro/clash/component/dialer"
	"github.com/igoogolx/itun2socks/internal/configuration"
	"github.com/igoogolx/itun2socks/pkg/log"
	tun "github.com/sagernet/sing-tun"
	"github.com/sagernet/sing/common/control"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/x/list"
	"github.com/sirupsen/logrus"
	"go.uber.org/atomic"
)

var defaultInterfaceName = atomic.NewString("")
var defaultInterfaceMonitor tun.DefaultInterfaceMonitor
var networkUpdateMonitor tun.NetworkUpdateMonitor

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

var monitorCallback *list.Element[tun.DefaultInterfaceUpdateCallback]

func StartMonitor() error {
	setting, err := configuration.GetSetting()
	if err != nil {
		return err
	}
	if len(setting.DefaultInterface) != 0 {
		update(setting.DefaultInterface)
		return nil
	}
	networkUpdateMonitor, err = tun.NewNetworkUpdateMonitor(logrus.StandardLogger())
	if err != nil {
		err = E.Cause(err, "create NetworkUpdateMonitor")
		return err
	}
	err = networkUpdateMonitor.Start()
	if err != nil {
		err = E.Cause(err, "start NetworkUpdateMonitor")
		return err
	}

	defaultInterfaceMonitor, err = tun.NewDefaultInterfaceMonitor(
		networkUpdateMonitor,
		logrus.StandardLogger(),
		tun.DefaultInterfaceMonitorOptions{
			OverrideAndroidVPN: true,
			InterfaceFinder:    control.NewDefaultInterfaceFinder(),
		})
	if err != nil {
		err = E.Cause(err, "create DefaultInterfaceMonitor")
		return err
	}
	monitorCallback = defaultInterfaceMonitor.RegisterCallback(func(defaultInterface *control.Interface, flags int) {
		//FIXME: flags?
		if defaultInterface != nil {
			update(defaultInterface.Name)
		}
	})
	err = defaultInterfaceMonitor.Start()
	if err != nil {
		return err
	}
	if defaultInterfaceMonitor.DefaultInterface() != nil {
		update(defaultInterfaceMonitor.DefaultInterface().Name)
	}
	return nil
}

func StopMonitor() error {
	defer func() {
		defaultInterfaceMonitor = nil
		networkUpdateMonitor = nil
		monitorCallback = nil
	}()
	if monitorCallback != nil {
		defaultInterfaceMonitor.UnregisterCallback(monitorCallback)
	}
	if networkUpdateMonitor != nil {
		err := networkUpdateMonitor.Close()
		if err != nil {
			return err
		}
	}
	if defaultInterfaceMonitor != nil {
		err := defaultInterfaceMonitor.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func update(name string) {
	defaultInterfaceName.Store(name)
	dialer.DefaultInterface.Store(name)
	log.Infoln(log.FormatLog(log.ExecutorPrefix, "update default interface: %v"), name)
}
