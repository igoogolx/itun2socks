package rule

import (
	"regexp"
	"strings"
)

type Domain struct {
	payload string
	policy  string
}

func (d Domain) Policy() string {
	return d.policy
}

func (d Domain) Match(value string) bool {
	return isContainsDomain(d.payload, value)
}

func NewDomainRule() {
}

func isContainsDomain(domain string, s string) bool {
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
