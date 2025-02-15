package conn

import (
	"fmt"
	"github.com/Dreamacro/clash/adapter"
	"github.com/Dreamacro/clash/adapter/outbound"
	C "github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/cfg/distribution/rule_engine"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/pkg/log"
	"strings"
	"sync"
)

var (
	proxies map[constants.Policy]C.Proxy
	mux     sync.RWMutex
)

type Matcher func(metadata *C.Metadata, rule rule_engine.Rule) (rule_engine.Rule, error)

func RejectQuicMather(metadata *C.Metadata, prevRule rule_engine.Rule) (rule_engine.Rule, error) {
	if prevRule.GetPolicy() == constants.PolicyProxy && strings.Contains(metadata.NetWork.String(), "udp") && metadata.DstPort.String() == "443" {
		log.Debugln("reject quic conn:%v", metadata.RemoteAddress())
		return rule_engine.BuiltInRejectRule, nil
	}
	return nil, fmt.Errorf("not quic")
}

func UpdateProxy(remoteProxy C.Proxy) {
	mux.Lock()
	defer mux.Unlock()
	proxies = make(map[constants.Policy]C.Proxy)
	proxies[constants.PolicyProxy] = remoteProxy
	proxies[constants.PolicyDirect] = adapter.NewProxy(outbound.NewDirect())
	proxies[constants.PolicyReject] = adapter.NewProxy(outbound.NewReject())
}

func GetProxy(rule constants.Policy) (C.Proxy, error) {
	mux.RLock()
	defer mux.RUnlock()
	connDialer := proxies[rule]
	if connDialer == nil {
		return nil, fmt.Errorf("empty dialer")
	}
	return connDialer, nil
}
