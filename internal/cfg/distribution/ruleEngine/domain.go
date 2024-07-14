package ruleEngine

import (
	"fmt"
	"github.com/igoogolx/itun2socks/internal/constants"
	"regexp"
	"strings"
)

type Domain struct {
	RuleType constants.RuleType `json:"ruleType"`
	Payload  string             `json:"payload"`
	Policy   constants.Policy   `json:"policy"`
}

func (d Domain) GetPolicy() constants.Policy {
	return d.Policy
}

func (d Domain) Type() constants.RuleType {
	return d.RuleType
}

func (d Domain) Match(value string) bool {
	return isContainsDomain(d.RuleType, d.Payload, value)
}

func (d Domain) Value() string {
	return d.Payload
}

func (d Domain) Valid() bool {
	return len(d.Payload) != 0
}

func NewDomainRule(ruleType constants.RuleType, payload string, policy constants.Policy) (*Domain, error) {
	rule := &Domain{ruleType, payload, policy}

	if !rule.Valid() {
		return nil, fmt.Errorf("invalid domain rule")
	}

	return rule, nil
}

func isContainsDomain(rType constants.RuleType, value string, s string) bool {
	switch rType {
	case constants.RuleDomainKeyword:
		return strings.Contains(s, value)
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
