package conn

import (
	"github.com/Dreamacro/clash/adapter"
	"github.com/Dreamacro/clash/adapter/outbound"
	C "github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/constants"
	"net/netip"
	"sync"
)

var (
	proxies map[constants.IpRule]C.Proxy
	mux     sync.RWMutex
)

const (
	ProxyAddrDomain ProxyAddrType = 0
	ProxyAddrIp     ProxyAddrType = 1
)

type ProxyAddrType int

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
	proxies = make(map[constants.IpRule]C.Proxy)
	proxies[constants.DistributionProxy] = proxy
	proxies[constants.DistributionBypass] = adapter.NewProxy(outbound.NewDirect())
}

func getProxy(rule constants.IpRule) C.Proxy {
	mux.RLock()
	defer mux.RUnlock()
	return proxies[rule]
}
