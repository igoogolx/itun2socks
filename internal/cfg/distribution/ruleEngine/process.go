package ruleEngine

import (
	"github.com/igoogolx/itun2socks/internal/constants"
	"path/filepath"
	"strings"
)

type Process struct {
	RuleType constants.RuleType `json:"ruleType"`
	Payload  string             `json:"payload"`
	Policy   constants.Policy   `json:"policy"`
}

func (p Process) GetPolicy() constants.Policy {
	return p.Policy
}

func (p Process) Type() constants.RuleType {
	return constants.RuleProcess
}

func (p Process) Match(value string) bool {
	return strings.EqualFold(filepath.Base(value), p.Payload)
}

func (p Process) Value() string {
	return p.Payload
}

func NewProcessRule(ruleType constants.RuleType, payload string, policy constants.Policy) (*Process, error) {
	return &Process{ruleType, payload, policy}, nil
}
