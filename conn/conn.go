package conn

import (
	C "github.com/Dreamacro/clash/constant"
	"net/netip"
	"sync"
)

var defaultProxy C.Proxy
var mux sync.RWMutex

type ProxyAddrType int

const (
	ProxyAddrDomain ProxyAddrType = 0
	ProxyAddrIp     ProxyAddrType = 1
)

type ProxyAddr struct {
	addr     string
	addrType ProxyAddrType
}

func NewProxyAddr(addr string) ProxyAddr {
	_, err := netip.ParseAddr(addr)
	var addrType ProxyAddrType
	if err == nil {
		addrType = ProxyAddrIp
	} else {
		addrType = ProxyAddrDomain
	}
	return ProxyAddr{addr, addrType}
}

func (p ProxyAddr) Addr() string {
	return p.addr
}

func (p ProxyAddr) Type() ProxyAddrType {
	return p.addrType
}

func UpdateProxy(proxy C.Proxy) {
	mux.Lock()
	defer mux.Unlock()
	defaultProxy = proxy
}

func getProxy() C.Proxy {
	mux.RLock()
	defer mux.RUnlock()
	return defaultProxy
}
