package dns

import (
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/igoogolx/itun2socks/internal/cfg/distribution/ruleEngine"
)

var dnsDomainCache, _ = lru.New[string, string](4 * 1024)
var dnsRuleCache, _ = lru.New[string, ruleEngine.Rule](4 * 1024)

func ResetCache() {
	dnsDomainCache.Purge()
	dnsRuleCache.Purge()
}

func GetCachedDnsItem(ip string) (string, ruleEngine.Rule, bool) {
	rawCachedDomain, ok := dnsDomainCache.Get(ip)
	if !ok {
		return "", nil, false
	}
	rawCachedRule, ok := dnsRuleCache.Get(ip)
	if !ok {
		return "", nil, false
	}
	return rawCachedDomain, rawCachedRule, true
}

func AddCachedDnsItem(ip, domain string, rule ruleEngine.Rule) {
	dnsDomainCache.Add(ip, domain)
	dnsRuleCache.Add(ip, rule)
}
