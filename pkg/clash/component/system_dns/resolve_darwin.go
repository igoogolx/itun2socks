package system_dns

import (
	"bufio"
	"context"
	"fmt"
	"net/netip"
	"os"
	"os/exec"
	"strings"

	"github.com/samber/lo"
)

func ResolveServers(_ string) ([]string, error) {

	var servers []string

	serviceServers, serviceErr := resolveServersFormService()

	if serviceErr == nil {
		servers = append(servers, serviceServers...)
	}

	resolvConfServers, resolvErr := dnsReadConfig()

	if resolvErr == nil {
		servers = append(servers, resolvConfServers...)
	}

	if len(servers) == 0 {
		return nil, fmt.Errorf("service error: %v, resolv conf error: %v", serviceErr, resolvErr)
	}

	return lo.Uniq(servers), nil

}

func resolveServersFormService() ([]string, error) {
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

//Copied from https://github.com/TransactCharlie/go-osx-dns

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

const resolvConf = "/etc/resolv.conf"

//Copied from https://github.com/MetaCubeX/mihomo/blob/Meta/dns/system_posix.go

func dnsReadConfig() (servers []string, err error) {
	file, err := os.Open(resolvConf)
	if err != nil {
		err = fmt.Errorf("failed to read %s: %w", resolvConf, err)
		return
	}
	defer func() { _ = file.Close() }()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && (line[0] == ';' || line[0] == '#') {
			// comment.
			continue
		}
		f := strings.Fields(line)
		if len(f) < 1 {
			continue
		}
		switch f[0] {
		case "nameserver": // add one name server
			if len(f) > 1 {
				if addr, err := netip.ParseAddr(f[1]); err == nil {
					servers = append(servers, addr.String())
				}
			}
		}
	}
	return
}
