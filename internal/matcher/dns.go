package matcher

import (
	"github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/constants"
	"sync"
)

type Dns interface {
	GetDnsType(question string, metadata *constant.Metadata) (constants.DnsType, error)
}

var defaultDnsMatcher Dns
var mux sync.RWMutex

func UpdateDnsMatcher(m Dns) {
	mux.Lock()
	defer mux.Unlock()
	defaultDnsMatcher = m
}

func GetDnsMatcher() Dns {
	mux.RLock()
	defer mux.RUnlock()
	return defaultDnsMatcher
}
