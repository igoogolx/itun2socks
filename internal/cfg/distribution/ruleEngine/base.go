package ruleEngine

import (
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"github.com/igoogolx/itun2socks/internal/constants"
)

type Rule interface {
	Match(value string) bool
	Value() string
	GetPolicy() constants.Policy
	Type() constants.RuleType
}

type Engine struct {
	rules []Rule
	cache *lru.Cache
}

func (e *Engine) Match(value string) (Rule, error) {
	cachedRule, ok := e.cache.Get(value)
	if ok {
		return cachedRule.(Rule), nil
	}
	for _, rule := range e.rules {
		if rule.Match(value) {
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
	cache, err := lru.New(1024)
	if err != nil {
		return nil, err
	}
	return &Engine{rules, cache}, nil
}
