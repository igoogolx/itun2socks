package configuration

import (
	"fmt"
	"github.com/igoogolx/itun2socks/internal/cfg/distribution/rule_engine"
	"slices"
	"strings"
)

func GetSelectedRule() (string, error) {
	c, err := Read()
	if err != nil {
		return "", err
	}
	return c.Selected.Rule, nil
}

func GetRuleIds() ([]string, error) {
	return rule_engine.GetRuleIds()
}

func AddCustomizedRule(rules []string) error {
	c, err := Read()
	if err != nil {
		return err
	}
	for _, rule := range rules {
		formatedRule := strings.TrimSpace(rule)
		targetIndex := slices.Index(c.Rules, formatedRule)
		if targetIndex != -1 {
			return fmt.Errorf("duplicated rule: %v", rule)
		}
		_, err = rule_engine.ParseRawValue(formatedRule)
		if err != nil {
			return err
		}
		c.Rules = append(c.Rules, formatedRule)
	}

	return Write(c)
}

func DeleteCustomizedRule(rules []string) error {
	c, err := Read()
	if err != nil {
		return err
	}
	for _, rule := range rules {
		_, err = rule_engine.ParseRawValue(rule)
		if err != nil {
			return err
		}
		var newRules []string
		for _, item := range c.Rules {
			if item != rule {
				newRules = append(newRules, item)
			}
		}
		c.Rules = newRules
	}
	return Write(c)
}

func GetCustomizedRules() ([]rule_engine.Rule, error) {
	c, err := Read()
	if err != nil {
		return nil, err
	}
	var items []rule_engine.Rule
	for _, rule := range c.Rules {
		item, err := rule_engine.ParseRawValue(rule)
		if err == nil {
			items = append(items, item)
		}
	}
	return items, nil
}

func EditCustomizedRule(oldRule string, newRule string) error {
	c, err := Read()
	if err != nil {
		return err
	}
	targetIndex := slices.Index(c.Rules, oldRule)
	if targetIndex != -1 {
		slices.Replace(c.Rules, targetIndex, targetIndex+1, newRule)
	}
	return Write(c)
}

func GetBuiltInRules(id string) ([]rule_engine.Rule, error) {
	return rule_engine.Parse(id, []string{})
}
