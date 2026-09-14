package conn

import (
	"sync"

	"github.com/igoogolx/itun2socks/internal/cfg/distribution/rule_engine"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/internal/dns"
	"github.com/igoogolx/itun2socks/pkg/log"
	metaC "github.com/metacubex/mihomo/constant"
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

func GetIsDNSConn(metadata *metaC.Metadata) bool {
	return metadata.NetWork == metaC.UDP && metadata.DstPort == constants.DnsPort
}

func resolveMetadata(metadata *metaC.Metadata) rule_engine.Rule {

	var logType = log.TcpPrefix
	var printLog = log.Infoln

	var rule rule_engine.Rule = rule_engine.BuiltInProxyRule
	for _, matcher := range GetConnMatcher() {
		tempRule, err := matcher(metadata, rule)
		if err == nil {
			rule = tempRule
		}
	}
	remoteAddr := metadata.DstIP.String()
	if metadata.NetWork == metaC.UDP {
		logType = log.UdpPrefix
	}
	if GetIsDNSConn(metadata) {
		printLog = log.Debugln
	}
	cachedDomain, ok := dns.GetCachedDnsItem(remoteAddr)
	if ok {
		printLog(log.FormatLog(logType, " %s to %s(%s) using %s"), metadata.SourceAddress(), metadata.RemoteAddress(), cachedDomain, rule.GetPolicy())
	} else {
		printLog(log.FormatLog(logType, " %s to %s using %s"), metadata.SourceAddress(), metadata.RemoteAddress(), rule.GetPolicy())
	}
	return rule
}
