package rule_engine

import (
	"fmt"
	"github.com/igoogolx/itun2socks/internal/constants"
	"os"
	"path/filepath"
	"runtime"
	"slices"
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
	processName := filepath.Base(value)

	if runtime.GOOS == "darwin" {
		parts := strings.Split(value, string(os.PathSeparator))
		targetIndex := slices.IndexFunc(parts, func(s string) bool {
			return strings.HasSuffix(s, ".app")
		})
		if targetIndex != -1 {
			processName = parts[targetIndex]
		}
	}

	return strings.EqualFold(processName, p.Payload)
}

func (p Process) Value() string {
	return p.Payload
}

func (p Process) Valid() bool {
	return len(p.Payload) != 0
}

func NewProcessRule(ruleType constants.RuleType, payload string, policy constants.Policy) (*Process, error) {
	rule := &Process{ruleType, strings.TrimSpace(payload), policy}
	if !rule.Valid() {
		return nil, fmt.Errorf("invalid process rule")
	}
	return rule, nil
}
