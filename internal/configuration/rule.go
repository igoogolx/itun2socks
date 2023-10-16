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

func SetRules(rules []string) error {
	c, err := Read()
	if err != nil {
		return err
	}
	return Write(c)
}

func GetRules() ([]string, error) {
	c, err := Read()
	if err != nil {
		return nil, err
	}
	return c.Rules, nil
}

func GetBuiltInRules(id string) ([]string, error) {
	c, err := Read()
	if err != nil {
		return nil, err
	}
	return c.Rules, nil
}
