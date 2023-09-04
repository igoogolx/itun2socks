package conn

import (
	"github.com/Dreamacro/clash/adapter"
	"github.com/Dreamacro/clash/adapter/outbound"
	C "github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/pkg/log"
	"go.uber.org/atomic"
	"net"
	"strings"
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

func getProxy(rule constants.IpRule) C.Proxy {
	mux.RLock()
	defer mux.RUnlock()

	if rule == constants.DistributionProxy {
		var remoteProxy = proxies[rule]
		selectedProxyAddr := remoteProxy.Addr()
		if len(selectedProxyAddr) == 0 {
			selectedProxyAddr = remoteProxy.Unwrap(&C.Metadata{}).Addr()
		}
		addr, _, err := net.SplitHostPort(selectedProxyAddr)
		if err == nil {
			log.Debugln(log.FormatLog(log.RulePrefix, "update proxy addr: %v"), addr)
			proxyAddr.Store(addr)
		} else {
			log.Errorln(log.FormatLog(log.RulePrefix, "invalid proxy addr: %v"), addr)
		}
	}

	return proxies[rule]
}

func GetProxyAddr() string {
	return proxyAddr.Load()
}

func GetIsProxyAddr(addr string) bool {
	storedAddr := proxyAddr.Load()
	if len(storedAddr) != 0 {
		return strings.Contains(storedAddr, addr)
	}
	return false
}
