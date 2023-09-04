package distribution

import (
	"github.com/igoogolx/itun2socks/pkg/log"
	"net/netip"
	"regexp"
	"strings"
)

type MatcherList interface {
	Has(s string) bool
}

func IsSubnetContainsIp(s string, i string) bool {
	prefix, err := netip.ParsePrefix(s)
	if err != nil {
		return s == i
	}
	ip, err := netip.ParseAddr(i)
	if err != nil {
		return s == i
	}
	return prefix.Contains(ip)
}

func IsDomainMatchRule(pattern string, domain string) bool {
	matched, err := regexp.MatchString(pattern, domain)
	if err != nil {
		log.Warnln(log.FormatLog(log.RulePrefix, "matching domain %s with pattern %s: %s"), domain, pattern, err)
	}
	return matched
}

func IsContainsDomain(domain string, s string) bool {
	i := strings.LastIndexByte(domain, '/')
	if i < 0 {
		return false
	}
	domainValue := domain[:i]
	domainType := domain[i+1:]
	switch domainType {
	case "0":
		return strings.Contains(domainValue, s)
	case "1":
		pattern, err := regexp.Compile(domainValue)
		if err != nil {
			return false
		}
		return pattern.MatchString(s)
	case "2":
		if !strings.HasSuffix(s, domainValue) {
			return false
		}
		return len(s) == len(domainValue) || s[len(s)-len(domainValue)-1] == '.'
	case "3":
		return domainValue == s
	}
	return false
}
