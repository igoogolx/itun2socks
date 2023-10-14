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

func GetRules() ([]string, error) {
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

func Parse(name string) ([]Rule, error) {
	items, err := readFile("rules/" + name)
	if err != nil {
		return nil, err
	}
	var rules []Rule
	for _, line := range items {
		chunks := trimArr(strings.Split(strings.TrimSpace(line), ","))
		if len(chunks) != 3 {
			break
		}
		ruleType := constants.RuleConfig(chunks[0])

		var rule Rule
		var err error
		switch ruleType {
		case constants.RuleIpCidr:
			rule, err = NewIpCidrRule(chunks[1], chunks[2])
			break
		case constants.RuleDomain:
			rule, err = NewDomainRule(chunks[1], chunks[2])
			break
		default:
			err = fmt.Errorf("rule type not match: %v", ruleType)

		}
		if err == nil {
			rules = append(rules, rule)
		}
	}
	return rules, nil
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
