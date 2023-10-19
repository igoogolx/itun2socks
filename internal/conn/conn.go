package conn

import (
	"github.com/Dreamacro/clash/adapter"
	"github.com/Dreamacro/clash/adapter/outbound"
	C "github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/constants"
	"sync"
)

var (
	proxies map[constants.RuleType]C.Proxy
	mux     sync.RWMutex
)

func UpdateProxy(remoteProxy C.Proxy) {
	mux.Lock()
	defer mux.Unlock()
	proxies = make(map[constants.RuleType]C.Proxy)
	proxies[constants.RuleProxy] = remoteProxy
	proxies[constants.RuleBypass] = adapter.NewProxy(outbound.NewDirect())
}

func GetProxy(rule constants.RuleType) C.Proxy {
	mux.RLock()
	defer mux.RUnlock()
	return proxies[rule]
}
