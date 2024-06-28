package ruleEngine

import (
	"fmt"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/igoogolx/itun2socks/internal/constants"
	"slices"
)

type Rule interface {
	Match(value string) bool
	Value() string
	GetPolicy() constants.Policy
	Type() constants.RuleType
}

type Engine struct {
	rules []Rule
	cache *lru.Cache[string, Rule]
}

func (e *Engine) Match(value string, types []constants.RuleType) (Rule, error) {
	cachedRule, ok := e.cache.Get(value)
	if ok {
		return cachedRule, nil
	}
	for _, rule := range e.rules {
		if slices.Contains(types, rule.Type()) && rule.Match(value) {
			e.cache.Add(value, rule)
			return rule, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func New(name string, extraRules []string) (*Engine, error) {
	rules, err := Parse(name, extraRules)
	if err != nil {
		return nil, err
	}
	cache, err := lru.New[string, Rule](1024)
	if err != nil {
		return nil, err
	}
	return &Engine{rules, cache}, nil
}
