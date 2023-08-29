package executor

import (
	"context"
	"github.com/igoogolx/itun2socks/cfg"
	network_iface "github.com/igoogolx/itun2socks/components/network-iface"
	"github.com/igoogolx/itun2socks/components/proxy-handler"
	"github.com/igoogolx/itun2socks/conn"
	"github.com/igoogolx/itun2socks/dns"
	localserver "github.com/igoogolx/itun2socks/local-server"
	"github.com/igoogolx/itun2socks/tunnel"
	sTun "github.com/sagernet/sing-tun"
	F "github.com/sagernet/sing/common/format"
	"net"
	"net/netip"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func CalculateInterfaceName(name string) (tunName string) {
	if runtime.GOOS == "darwin" {
		tunName = "utun"
	} else if name != "" {
		tunName = name
		return
	} else {
		tunName = "tun"
	}
	interfaces, err := net.Interfaces()
	if err != nil {
		return
	}
	var tunIndex int
	for _, netInterface := range interfaces {
		if strings.HasPrefix(netInterface.Name, tunName) {
			index, parseErr := strconv.ParseInt(netInterface.Name[len(tunName):], 10, 16)
			if parseErr == nil {
				tunIndex = int(index) + 10
			}
		}
	}
	tunName = F.ToString(tunName, tunIndex)
	return
}

func New() (*Client, error) {
	config, err := cfg.New()
	if err != nil {
		return nil, err
	}
	tunDevice := config.Device
	tunInterfaceName := CalculateInterfaceName(tunDevice.Name)
	tunOptions := sTun.Options{
		Name:         tunInterfaceName,
		MTU:          uint32(tunDevice.Mtu),
		Inet4Address: []netip.Prefix{netip.MustParsePrefix(tunDevice.Gateway.String())},
		AutoRoute:    true,
		StrictRoute:  true,
	}
	tun, err := sTun.New(tunOptions)
	if err != nil {
		return nil, err
	}
	stack, err := sTun.NewStack("gvisor", sTun.StackOptions{
		Context:    context.TODO(),
		Handler:    proxy_handler.New(tunnel.TcpQueue(), tunnel.UdpQueue()),
		Tun:        tun,
		Name:       tunInterfaceName,
		MTU:        uint32(tunDevice.Mtu),
		UDPTimeout: int64(5 * time.Minute),
	})
	if err != nil {
		return nil, err
	}

	newLocalServer := localserver.New(config.LocalServer.HttpAddr)
	interfaceHandler, err := network_iface.New()
	if err != nil {
		return nil, err
	}
	err = interfaceHandler.Monitor.Start()
	if err != nil {
		return nil, err
	}
	if err = updateCfg(config); err != nil {
		return nil, err
	}
	return &Client{
		stack:                   stack,
		tun:                     tun,
		localserver:             newLocalServer,
		defaultInterfaceHandler: *interfaceHandler,
		deviceName:              tunInterfaceName,
		localDns:                config.Rule.Dns.Local.Client.Nameservers(),
		remoteDns:               config.Rule.Dns.Remote.Client.Nameservers(),
		boostDns:                config.Rule.Dns.BoostNameserver,
	}, nil
}

func updateCfg(config cfg.Config) error {
	updateMatcher(config)
	updateConn(config)
	updateDns(config)
	return nil
}

func updateMatcher(c cfg.Config) {
	conn.UpdateMatcher(c.Rule)
}

func updateDns(c cfg.Config) {
	dns.UpdateMatcher(c.Rule)
}

func updateConn(c cfg.Config) {
	conn.UpdateProxy(c.Proxy)
}
