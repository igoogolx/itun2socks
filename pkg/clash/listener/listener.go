package listener

import (
	"fmt"
	"net"
	"strconv"
	"sync"

	"github.com/igoogolx/itun2socks/pkg/clash/adapter/inbound"
	C "github.com/igoogolx/itun2socks/pkg/clash/constant"
	"github.com/igoogolx/itun2socks/pkg/clash/listener/http"
	"github.com/igoogolx/itun2socks/pkg/clash/listener/mixed"
	"github.com/igoogolx/itun2socks/pkg/clash/listener/redir"
	"github.com/igoogolx/itun2socks/pkg/clash/listener/socks"
	"github.com/igoogolx/itun2socks/pkg/clash/listener/tproxy"
	"github.com/igoogolx/itun2socks/pkg/clash/listener/tunnel"
	"github.com/igoogolx/itun2socks/pkg/clash/log"

	"github.com/samber/lo"
)

var (
	allowLan    = false
	bindAddress = "*"

	tcpListeners = map[C.Inbound]C.Listener{}
	udpListeners = map[C.Inbound]C.Listener{}

	tunnelTCPListeners = map[string]*tunnel.Listener{}
	tunnelUDPListeners = map[string]*tunnel.PacketConn{}

	// lock for recreate function
	recreateMux sync.Mutex
	tunnelMux   sync.Mutex
)

type Ports struct {
	Port       int `json:"port"`
	SocksPort  int `json:"socks-port"`
	RedirPort  int `json:"redir-port"`
	TProxyPort int `json:"tproxy-port"`
	MixedPort  int `json:"mixed-port"`
}

var tcpListenerCreators = map[C.InboundType]tcpListenerCreator{
	C.InboundTypeHTTP:   http.New,
	C.InboundTypeSocks:  socks.New,
	C.InboundTypeRedir:  redir.New,
	C.InboundTypeTproxy: tproxy.New,
	C.InboundTypeMixed:  mixed.New,
}

var udpListenerCreators = map[C.InboundType]udpListenerCreator{
	C.InboundTypeSocks:  socks.NewUDP,
	C.InboundTypeRedir:  tproxy.NewUDP,
	C.InboundTypeTproxy: tproxy.NewUDP,
	C.InboundTypeMixed:  socks.NewUDP,
}

type (
	tcpListenerCreator func(addr string, tcpIn chan<- C.ConnContext) (C.Listener, error)
	udpListenerCreator func(addr string, udpIn chan<- *inbound.PacketAdapter) (C.Listener, error)
)

func AllowLan() bool {
	return allowLan
}

func BindAddress() string {
	return bindAddress
}

func SetAllowLan(al bool) {
	allowLan = al
}

func SetBindAddress(host string) {
	bindAddress = host
}

func createListener(inbound C.Inbound, tcpIn chan<- C.ConnContext, udpIn chan<- *inbound.PacketAdapter) {
	addr := inbound.BindAddress
	if portIsZero(addr) {
		return
	}
	tcpCreator := tcpListenerCreators[inbound.Type]
	udpCreator := udpListenerCreators[inbound.Type]
	if tcpCreator == nil && udpCreator == nil {
		log.Errorln("inbound type %s is not supported", inbound.Type)
		return
	}
	if tcpCreator != nil {
		tcpListener, err := tcpCreator(addr, tcpIn)
		if err != nil {
			log.Errorln("create addr %s tcp listener error: %v", addr, err)
			return
		}
		tcpListeners[inbound] = tcpListener
	}
	if udpCreator != nil {
		udpListener, err := udpCreator(addr, udpIn)
		if err != nil {
			log.Errorln("create addr %s udp listener error: %v", addr, err)
			return
		}
		udpListeners[inbound] = udpListener
	}
	log.Infoln("inbound %s created successfully", inbound.ToAlias())
}

func closeListener(inbound C.Inbound) {
	listener := tcpListeners[inbound]
	if listener != nil {
		if err := listener.Close(); err != nil {
			log.Errorln("close tcp address `%s` error: %s", inbound.ToAlias(), err.Error())
		}
		delete(tcpListeners, inbound)
	}
	listener = udpListeners[inbound]
	if listener != nil {
		if err := listener.Close(); err != nil {
			log.Errorln("close udp address `%s` error: %s", inbound.ToAlias(), err.Error())
		}
		delete(udpListeners, inbound)
	}
}

func getNeedCloseAndCreateInbound(originInbounds []C.Inbound, newInbounds []C.Inbound) ([]C.Inbound, []C.Inbound) {
	needCloseMap := map[C.Inbound]bool{}
	needClose := []C.Inbound{}
	needCreate := []C.Inbound{}

	for _, inbound := range originInbounds {
		needCloseMap[inbound] = true
	}
	for _, inbound := range newInbounds {
		if needCloseMap[inbound] {
			delete(needCloseMap, inbound)
		} else {
			needCreate = append(needCreate, inbound)
		}
	}
	for inbound := range needCloseMap {
		needClose = append(needClose, inbound)
	}
	return needClose, needCreate
}

// only recreate inbound config listener
func ReCreateListeners(inbounds []C.Inbound, tcpIn chan<- C.ConnContext, udpIn chan<- *inbound.PacketAdapter) {
	newInbounds := []C.Inbound{}
	newInbounds = append(newInbounds, inbounds...)
	for _, inbound := range getInbounds() {
		if inbound.IsFromPortCfg {
			newInbounds = append(newInbounds, inbound)
		}
	}
	reCreateListeners(newInbounds, tcpIn, udpIn)
}

// only recreate ports config listener
func ReCreatePortsListeners(ports Ports, tcpIn chan<- C.ConnContext, udpIn chan<- *inbound.PacketAdapter) {
	newInbounds := []C.Inbound{}
	newInbounds = append(newInbounds, GetInbounds()...)
	newInbounds = addPortInbound(newInbounds, C.InboundTypeHTTP, ports.Port)
	newInbounds = addPortInbound(newInbounds, C.InboundTypeSocks, ports.SocksPort)
	newInbounds = addPortInbound(newInbounds, C.InboundTypeRedir, ports.RedirPort)
	newInbounds = addPortInbound(newInbounds, C.InboundTypeTproxy, ports.TProxyPort)
	newInbounds = addPortInbound(newInbounds, C.InboundTypeMixed, ports.MixedPort)
	reCreateListeners(newInbounds, tcpIn, udpIn)
}

func addPortInbound(inbounds []C.Inbound, inboundType C.InboundType, port int) []C.Inbound {
	if port != 0 {
		inbounds = append(inbounds, C.Inbound{
			Type:          inboundType,
			BindAddress:   genAddr(bindAddress, port, allowLan),
			IsFromPortCfg: true,
		})
	}
	return inbounds
}

func reCreateListeners(inbounds []C.Inbound, tcpIn chan<- C.ConnContext, udpIn chan<- *inbound.PacketAdapter) {
	recreateMux.Lock()
	defer recreateMux.Unlock()
	needClose, needCreate := getNeedCloseAndCreateInbound(getInbounds(), inbounds)
	for _, inbound := range needClose {
		closeListener(inbound)
	}
	for _, inbound := range needCreate {
		createListener(inbound, tcpIn, udpIn)
	}
}

func GetInbounds() []C.Inbound {
	return lo.Filter(getInbounds(), func(inbound C.Inbound, idx int) bool {
		return !inbound.IsFromPortCfg
	})
}

// GetInbounds return the inbounds of proxy servers
func getInbounds() []C.Inbound {
	var inbounds []C.Inbound
	for inbound := range tcpListeners {
		inbounds = append(inbounds, inbound)
	}
	for inbound := range udpListeners {
		if _, ok := tcpListeners[inbound]; !ok {
			inbounds = append(inbounds, inbound)
		}
	}
	return inbounds
}

// GetPorts return the ports of proxy servers
func GetPorts() *Ports {
	ports := &Ports{}
	for _, inbound := range getInbounds() {
		fillPort(inbound, ports)
	}
	return ports
}

func fillPort(inbound C.Inbound, ports *Ports) {
	if inbound.IsFromPortCfg {
		port := getPort(inbound.BindAddress)
		switch inbound.Type {
		case C.InboundTypeHTTP:
			ports.Port = port
		case C.InboundTypeSocks:
			ports.SocksPort = port
		case C.InboundTypeTproxy:
			ports.TProxyPort = port
		case C.InboundTypeRedir:
			ports.RedirPort = port
		case C.InboundTypeMixed:
			ports.MixedPort = port
		default:
			// do nothing
		}
	}
}

func portIsZero(addr string) bool {
	_, port, err := net.SplitHostPort(addr)
	if port == "0" || port == "" || err != nil {
		return true
	}
	return false
}

func genAddr(host string, port int, allowLan bool) string {
	if allowLan {
		if host == "*" {
			return fmt.Sprintf(":%d", port)
		}
		return fmt.Sprintf("%s:%d", host, port)
	}

	return fmt.Sprintf("127.0.0.1:%d", port)
}

func getPort(addr string) int {
	_, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return 0
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0
	}
	return port
}
