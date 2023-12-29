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

func (d Domain) GetPolicy() constants.Policy {
	return constants.Policy(d.Policy)
}

func (d Domain) Type() constants.RuleConfig {
	return constants.RuleDomain
}

func (d Domain) Match(value string) bool {
	return isContainsDomain(d.RuleType, d.Payload, value)
}

func (d Domain) Value() string {
	return d.Payload
}

func NewDomainRule(payload, policy string) (*Domain, error) {
	return &Domain{constants.RuleDomain, payload, policy}, nil
}

func isContainsDomain(rType constants.RuleConfig, value string, s string) bool {
	switch rType {
	case constants.RuleDomainKeyword:
		return strings.Contains(value, s)
	case constants.RuleDomainRegex:
		pattern, err := regexp.Compile(value)
		if err != nil {
			return false
		}
		return pattern.MatchString(s)
	case constants.RuleDomainSuffix:
		if !strings.HasSuffix(s, value) {
			return false
		}
		return len(s) == len(value) || s[len(s)-len(value)-1] == '.'
	case constants.RuleDomain:
		return value == s
	}
	return false
}
