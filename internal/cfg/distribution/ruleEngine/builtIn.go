package ruleEngine

import (
	"fmt"
	"github.com/igoogolx/itun2socks/internal/constants"
	"path/filepath"
	"strings"
)

type builtIn struct {
	RuleType constants.RuleType `json:"ruleType"`
	Payload  string             `json:"payload"`
	Policy   constants.Policy   `json:"policy"`
}

func (p builtIn) GetPolicy() constants.Policy {
	return p.Policy
}

func (p builtIn) Type() constants.RuleType {
	return p.RuleType
}

func (p builtIn) Match(value string) bool {
	return strings.EqualFold(filepath.Base(value), p.Payload)
}

func (p builtIn) Value() string {
	return p.Payload
}

func (p builtIn) Valid() bool {
	return true
}

func newBuiltIn(payload string, policy constants.Policy) (*builtIn, error) {
	rule := &builtIn{constants.RuleBuiltIn, payload, policy}
	if !rule.Valid() {
		return nil, fmt.Errorf("invalid builtin rule")
	}
	return rule, nil
}

var (
	BuiltInProxyRule, _  = newBuiltIn("*", constants.PolicyProxy)
	BuiltInRejectRule, _ = newBuiltIn("*", constants.PolicyReject)
)
