//go:build darwin

package dns

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
)

var originalDnsServers []string

func getOriginalDnsServers() ([]string, error) {
	var dnsServers []string
	cmd := exec.Command("networksetup", "-getdnsservers", getNetworkService())
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

func Hijack() error {
	var err error
	originalDnsServers, err = getOriginalDnsServers()
	if err != nil {
		return err
	}
	dnsServer := "8.8.8.8"
	cmd := exec.Command("networksetup", "-setdnsservers", getNetworkService(), dnsServer)
	_, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to hajack DNS servers: %v", err)
	}
	return nil
}

func Resume() error {
	var err error
	defer func() {
		originalDnsServers = []string{}
	}()
	dnsServer := "Empty"
	if len(originalDnsServers) != 0 {
		dnsServer = strings.Join(originalDnsServers, " ")
	}
	cmd := exec.Command("networksetup", "-setdnsservers", getNetworkService(), dnsServer)
	_, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to resume DNS servers: %v", err)
	}
	return nil
}
