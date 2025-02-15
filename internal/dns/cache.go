package dns

import (
	lru "github.com/hashicorp/golang-lru/v2"
)

var dnsDomainCache, _ = lru.New[string, string](4 * 1024)

func ResetCache() {
	dnsDomainCache.Purge()
}

func GetCachedDnsItem(ip string) (string, bool) {
	rawCachedDomain, ok := dnsDomainCache.Get(ip)
	if !ok {
		return "", false
	}
	return rawCachedDomain, true
}

func addCachedDnsItem(ip, domain string) {
	dnsDomainCache.Add(ip, domain)
}
