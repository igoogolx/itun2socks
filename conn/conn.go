package conn

import (
	C "github.com/Dreamacro/clash/constant"
	"sync"
)

var defaultProxy C.Proxy
var mux sync.RWMutex

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
