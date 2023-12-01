package conn

import (
	"fmt"
	"github.com/Dreamacro/clash/adapter"
	"github.com/Dreamacro/clash/adapter/outbound"
	C "github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/pkg/log"
	"strings"
	"sync"
)

var (
	proxies map[constants.RuleType]C.Proxy
	mux     sync.RWMutex
)

type Matcher func(metadata *C.Metadata) (constants.RuleType, error)

func RejectQuicMather(metadata *C.Metadata) (constants.RuleType, error) {
	if strings.Contains(metadata.NetWork.String(), "udp") && metadata.DstPort.String() == "443" {
		log.Debugln("reject quic conn:%v", metadata.RemoteAddress())
		return constants.RuleReject, nil
	}
	return constants.RuleProxy, fmt.Errorf("not quic")
}

func UpdateProxy(remoteProxy C.Proxy) {
	mux.Lock()
	defer mux.Unlock()
	proxies = make(map[constants.RuleType]C.Proxy)
	proxies[constants.RuleProxy] = remoteProxy
	proxies[constants.RuleBypass] = adapter.NewProxy(outbound.NewDirect())
	proxies[constants.RuleReject] = adapter.NewProxy(outbound.NewReject())
}

func GetProxy(rule constants.RuleType) C.Proxy {
	mux.RLock()
	defer mux.RUnlock()
	return proxies[rule]
}
