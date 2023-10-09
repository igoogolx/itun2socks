package rule

import (
	"fmt"
	"github.com/igoogolx/itun2socks/internal/constants"
	"strings"
)

func trimArr(arr []string) (r []string) {
	for _, e := range arr {
		r = append(r, strings.TrimSpace(e))
	}
	return
}

func Parse(lines []string) []Rule {
	var rules []Rule
	for _, line := range lines {
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
			rule, err = NewIpCidrRule(chunks[1], chunks[2])
			break
		default:
			err = fmt.Errorf("rule type not match: %v", ruleType)

		}
		if err == nil {
			rules = append(rules, rule)
		}
	}
	return rules
}
