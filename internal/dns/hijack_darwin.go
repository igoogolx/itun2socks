//go:build darwin

package dns

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
)

var originalDnsServers []string

func getOriginalDnsServers(networkService string) ([]string, error) {
	var dnsServers []string
	cmd := exec.Command("networksetup", "-getdnsservers", networkService)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve original DNS servers: %v", err)
	}
	rawDnsServers := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, server := range rawDnsServers {
		ip := net.ParseIP(server)
		if ip != nil {
			dnsServers = append(dnsServers, server)
		}
	}
	return dnsServers, nil
}

func Hijack(networkService string, dnsServer string) ([]string, error) {
	var err error
	originalDnsServers, err = getOriginalDnsServers(networkService)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command("networksetup", "-setdnsservers", networkService, dnsServer)
	_, err = cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to hajack DNS servers: %v", err)
	}
	return originalDnsServers, nil
}

func Resume(networkService string) error {
	var err error
	defer func() {
		originalDnsServers = []string{}
	}()
	dnsServer := "empty"
	if len(originalDnsServers) != 0 {
		dnsServer = strings.Join(originalDnsServers, " ")
	}
	cmd := exec.Command("networksetup", "-setdnsservers", networkService, dnsServer)
	_, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to resume DNS servers: %v", err)
	}
	return nil
}
