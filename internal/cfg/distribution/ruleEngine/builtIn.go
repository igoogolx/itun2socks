package ruleEngine

import (
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

func newBuiltIn(payload string, policy constants.Policy) (*builtIn, error) {
	return &builtIn{constants.RuleBuiltIn, payload, policy}, nil
}

var (
	BuiltInProxyRule, _  = newBuiltIn("*", constants.PolicyProxy)
	BuiltInDirectRule, _ = newBuiltIn("*", constants.PolicyDirect)
	BuiltInRejectRule, _ = newBuiltIn("*", constants.PolicyReject)
)
