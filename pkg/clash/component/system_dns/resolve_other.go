//go:build !windows && !darwin

package system_dns

import "fmt"

func ResolveServers(ifaceName string) ([]string, error) {
	return nil, fmt.Errorf("Not support resolve system dns servers!")
}
