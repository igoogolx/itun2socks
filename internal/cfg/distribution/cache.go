package distribution

import (
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/igoogolx/itun2socks/internal/constants"
)

var dnsDomainCache, _ = lru.New[string, string](4 * 1024)
var dnsRuleCache, _ = lru.New[string, constants.DnsType](4 * 1024)

func ResetCache() {
	dnsDomainCache.Purge()
	dnsRuleCache.Purge()
}

func GetCachedDnsItem(ip string) (string, constants.DnsType, bool) {
	rawCachedDomain, ok := dnsDomainCache.Get(ip)
	if !ok {
		return "", constants.LocalDns, false
	}
	rawCachedRule, ok := dnsRuleCache.Get(ip)
	if !ok {
		return "", constants.LocalDns, false
	}
	return rawCachedDomain, rawCachedRule, true
}

func AddCachedDnsItem(ip, domain string, rule constants.DnsType) {
	dnsDomainCache.Add(ip, domain)
	dnsRuleCache.Add(ip, rule)
}
