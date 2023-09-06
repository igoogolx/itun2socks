package conn

import (
	"github.com/Dreamacro/clash/adapter"
	"github.com/Dreamacro/clash/adapter/outbound"
	C "github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/constants"
	"go.uber.org/atomic"
	"sync"
)

var (
	proxies   map[constants.IpRule]C.Proxy
	mux       sync.RWMutex
	proxyAddr atomic.String
)

type ProxyAddrType int

type ProxyAddr struct {
	addr     string
	addrType ProxyAddrType
}

func (p ProxyAddr) Addr() string {
	return p.addr
}

func (p ProxyAddr) Type() ProxyAddrType {
	return p.addrType
}

func UpdateProxy(remoteProxy C.Proxy) {
	mux.Lock()
	defer mux.Unlock()
	proxies = make(map[constants.IpRule]C.Proxy)
	proxies[constants.DistributionProxy] = remoteProxy
	proxies[constants.DistributionBypass] = adapter.NewProxy(outbound.NewDirect())
}

func GetProxy(rule constants.IpRule) C.Proxy {
	mux.RLock()
	defer mux.RUnlock()
	return proxies[rule]
}

func GetProxyAddr() string {
	return proxyAddr.Load()
}
