package system_dns

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

func ResolveServers(_ string) ([]string, error) {
	service, err := getNetworkService()
	if err != nil {
		return nil, err
	}

	dnsPayload, err := getDNSForPrimaryService(service)
	if err != nil {
		return nil, err
	}
	return dnsPayload.ServerAddresses, nil
}

func execSCutilScript(ctx context.Context, script []string) ([]string, error) {
	cmd := exec.CommandContext(ctx, "scutil")

	stdIn, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	for _, line := range script {
		cmdLine := []byte(strings.TrimSpace(line) + "\n")
		n, err := stdIn.Write([]byte(strings.TrimSpace(line) + "\n"))
		if n != len(cmdLine) {
			return nil, fmt.Errorf("no all bytes written to scutil")
		}
		if err != nil {
			return nil, err
		}
	}

	buff := bufio.NewScanner(stdOut)
	var allText []string

	for buff.Scan() {
		allText = append(allText, buff.Text())
	}

	return allText, nil
}

func getNetworkService() (string, error) {
	script := []string{
		"open",
		"get State:/Network/Global/IPv4",
		"d.show",
		"close",
		"quit",
	}

	ipv4Settings, err := execSCutilScript(context.Background(), script)
	if err != nil {
		return "", err
	}

	for _, l := range ipv4Settings {
		parts := strings.Split(l, ":")
		if strings.TrimSpace(parts[0]) == "PrimaryService" {
			return strings.TrimSpace(parts[1]), nil
		}
	}

	return "", fmt.Errorf("didn't find primary service")
}

type DNSPayload struct {
	DomainName      string
	ServerAddresses []string
}

func getDNSForPrimaryService(service string) (*DNSPayload, error) {
	script := []string{
		"open",
		fmt.Sprintf("get State:/Network/Service/%s/DNS", service),
		"d.show",
		"close",
		"quit",
	}

	dnsSettings, err := execSCutilScript(context.Background(), script)
	if err != nil {
		return nil, err
	}

	var addresses []string
	domain := ""
	for i := 1; i < len(dnsSettings); i++ {
		parts := strings.Split(dnsSettings[i], ":")
		if strings.TrimSpace(parts[0]) == "ServerAddresses" {
			i += 1
			for strings.TrimSpace(dnsSettings[i]) != "}" {
				dnsPart := strings.Split(dnsSettings[i], ":")
				addresses = append(addresses, strings.TrimSpace(dnsPart[1]))
				i++
			}
		}
		if strings.TrimSpace(parts[0]) == "DomainName" {
			domain = strings.TrimSpace(parts[1])
		}
	}

	return &DNSPayload{
		DomainName:      domain,
		ServerAddresses: addresses,
	}, nil
}
