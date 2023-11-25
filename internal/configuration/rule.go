package configuration

import "github.com/igoogolx/itun2socks/internal/cfg/distribution/ruleEngine"

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

func AddCustomizedRule(rule string) error {
	c, err := Read()
	if err != nil {
		return err
	}
	_, err = ruleEngine.ParseRawValue(rule)
	if err != nil {
		return err
	}
	c.Rules = append(c.Rules, rule)
	return Write(c)
}

func DeleteCustomizedRule(rule string) error {
	c, err := Read()
	if err != nil {
		return err
	}
	_, err = ruleEngine.ParseRawValue(rule)
	if err != nil {
		return err
	}
	var rules []string
	for _, item := range c.Rules {
		if item != rule {
			rules = append(rules, item)
		}
	}
	c.Rules = rules
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

func GetBuiltInRules(id string) ([]ruleEngine.Rule, error) {
	return ruleEngine.Parse(id, []string{})
}
