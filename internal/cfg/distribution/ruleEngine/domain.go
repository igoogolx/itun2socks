package ruleEngine

import (
	"github.com/igoogolx/itun2socks/internal/constants"
	"regexp"
	"strings"
)

type Domain struct {
	RuleType constants.RuleConfig `json:"ruleType"`
	Payload  string               `json:"payload"`
	Policy   string               `json:"policy"`
}

func (d Domain) GetPolicy() constants.RuleType {
	return constants.RuleType(d.Policy)
}

func (d Domain) Type() constants.RuleConfig {
	return constants.RuleDomain
}

func (d Domain) Match(value string) bool {
	return isContainsDomain(d.Payload, value)
}

func (d Domain) Value() string {
	return d.Payload
}

func NewDomainRule(payload, policy string) (*Domain, error) {
	return &Domain{constants.RuleDomain, payload, policy}, nil
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
