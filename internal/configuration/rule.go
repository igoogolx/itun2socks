package configuration

import "github.com/igoogolx/itun2socks/internal/cfg/distribution/ruleEngine"

func GetSelectedRule() (string, error) {
	c, err := Read()
	if err != nil {
		return "", err
	}
	return c.Selected.Rule, nil
}

func GetRules() ([]string, error) {
	return ruleEngine.GetRules()
}
