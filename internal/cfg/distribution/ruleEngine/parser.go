package ruleEngine

import (
	"embed"
	"fmt"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/pkg/list"
	"github.com/igoogolx/itun2socks/pkg/log"
	"io/fs"
	"strings"
)

//go:embed rules/*
var data embed.FS

func trimArr(arr []string) (r []string) {
	for _, e := range arr {
		r = append(r, strings.TrimSpace(e))
	}
	return
}

func GetRuleIds() ([]string, error) {
	ruleFiles, err := data.ReadDir("rules")
	var rules []string
	if err != nil {
		return nil, err
	}
	for _, rule := range ruleFiles {
		rules = append(rules, rule.Name())
	}
	return rules, err
}

func Parse(name string, extraRules []string) ([]Rule, error) {
	var err error
	var rules []Rule
	builtInItems, err := readFile("rules/" + name)
	if err != nil {
		return nil, err
	}
	for _, line := range extraRules {
		rule, err := ParseRawValue(line)
		if err == nil {
			rules = append(rules, rule)
		}
	}
	for _, line := range builtInItems {
		rule, err := ParseRawValue(line)
		if err == nil {
			rules = append(rules, rule)
		}
	}
	return rules, nil
}

func ParseRawValue(line string) (Rule, error) {
	chunks := trimArr(strings.Split(strings.TrimSpace(line), ","))
	if len(chunks) != 3 {
		return nil, fmt.Errorf("invald rule line")
	}
	return ParseItem(chunks[0], chunks[1], chunks[2])

}

func ParseItem(rawRuleType, value, policy string) (Rule, error) {
	ruleType := constants.RuleType(rawRuleType)

	var rule Rule
	var err error
	switch ruleType {
	case constants.RuleIpCidr:
		rule, err = NewIpCidrRule(value, policy)
		break
	case constants.RuleDomain:
	case constants.RuleDomainKeyword:
	case constants.RuleDomainSuffix:
	case constants.RuleDomainRegex:
		rule, err = NewDomainRule(value, policy)
		break
	default:
		err = fmt.Errorf("rule type not match: %v", ruleType)
	}
	return rule, err
}

func readFile(path string) ([]string, error) {
	file, err := data.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file fs.File) {
		err := file.Close()
		if err != nil {
			log.Warnln(log.FormatLog(log.ConfigurationPrefix, "fail to close geo file: %v"), path)
		}
	}(file)
	items, err := list.ParseFile(file)
	return items, nil
}
