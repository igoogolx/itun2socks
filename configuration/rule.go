package configuration

import (
	"fmt"
	"github.com/gofrs/uuid"
)

func GetSelectedRule() (RuleCfg, error) {
	c, err := Read()
	if err != nil {
		return RuleCfg{}, err
	}
	for _, v := range c.Rule {
		if v.Id == c.Selected.Rule {
			return v, nil
		}
	}

	return RuleCfg{}, fmt.Errorf("error getting selected rule,id:%v,err:%v", c.Selected.Rule, err)
}

func GetRules() ([]RuleCfg, error) {
	c, err := Read()
	if err != nil {
		return nil, err
	}
	return c.Rule, nil
}

func DeleteRule(id string) error {
	c, err := Read()
	if err != nil {
		return err
	}
	rules := make([]RuleCfg, 0)
	for _, v := range c.Rule {
		if v.Id != id {
			rules = append(rules, v)
		}
	}
	c.Rule = rules
	err = Write(c)
	if err != nil {
		return err
	}
	return nil
}

func UpdateRule(id string, rule RuleCfg) error {
	c, err := Read()
	if err != nil {
		return err
	}
	for i, v := range c.Rule {
		if v.Id == id {
			c.Rule[i] = rule
			break
		}
	}
	err = Write(c)
	if err != nil {
		return err
	}
	return nil
}

func AddRule(rule RuleCfg) (string, error) {
	c, err := Read()
	if err != nil {
		return "", err
	}
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	rule.Id = id.String()
	c.Rule = append(c.Rule, rule)
	err = Write(c)
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
