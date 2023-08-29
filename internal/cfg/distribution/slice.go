package distribution

import "strings"

func IsDomainsContain(items []string, value string) bool {
	for _, item := range items {
		if strings.Contains(item, value) {
			return true
		}
	}
	return false
}
