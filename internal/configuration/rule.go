package configuration

import (
	"fmt"
	"github.com/igoogolx/itun2socks/internal/cfg/distribution/ruleEngine"
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
	return ruleEngine.GetRuleIds()
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
		_, err = ruleEngine.ParseRawValue(formatedRule)
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
		_, err = ruleEngine.ParseRawValue(rule)
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

func GetCustomizedRules() ([]ruleEngine.Rule, error) {
	c, err := Read()
	if err != nil {
		return nil, err
	}
	var items []ruleEngine.Rule
	for _, rule := range c.Rules {
		item, err := ruleEngine.ParseRawValue(rule)
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

func GetBuiltInRules(id string) ([]ruleEngine.Rule, error) {
	return ruleEngine.Parse(id, []string{})
}
