package conn

import (
	C "github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/cfg/distribution/ruleEngine"
	"github.com/igoogolx/itun2socks/internal/dns"
	"github.com/igoogolx/itun2socks/pkg/log"
	"sync"
)

var defaultConnMatchers []Matcher
var matcherMux sync.RWMutex

func UpdateConnMatcher(matchers []Matcher) {
	matcherMux.Lock()
	defer matcherMux.Unlock()
	defaultConnMatchers = matchers
}

func GetConnMatcher() []Matcher {
	matcherMux.RLock()
	defer matcherMux.RUnlock()
	return defaultConnMatchers
}

func resolveMetadata(metadata *C.Metadata) ruleEngine.Rule {
	var rule ruleEngine.Rule = ruleEngine.BuiltInProxyRule
	for _, matcher := range GetConnMatcher() {
		tempRule, err := matcher(metadata, rule)
		if err == nil {
			rule = tempRule
		}
	}
	var logType = log.TunPrefix
	if metadata.NetWork == C.UDP {
		logType = log.UdpPrefix
	}
	remoteAddr := metadata.DstIP.String()
	cachedDomain, _, ok := dns.GetCachedDnsItem(remoteAddr)
	if ok {
		log.Infoln(log.FormatLog(logType, " %s --> %s(%s) using %s"), metadata.SourceAddress(), metadata.RemoteAddress(), cachedDomain, rule.GetPolicy())
	} else {
		log.Infoln(log.FormatLog(logType, " %s --> %s using %s"), metadata.SourceAddress(), metadata.RemoteAddress(), rule.GetPolicy())
	}
	return rule
}
