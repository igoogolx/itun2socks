package distribution

import (
	"fmt"
	"github.com/Dreamacro/clash/component/fakeip"
	C "github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/cfg/distribution/rule_engine"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/internal/dns"
	"github.com/igoogolx/itun2socks/internal/matcher"
)

type Config struct {
	Dns DnsDistribution
}

func (c Config) ConnMatcher(metadata *C.Metadata, _ rule_engine.Rule) (rule_engine.Rule, error) {
	processPath := metadata.ProcessPath
	if len(processPath) != 0 {
		if rule, err := matcher.GetRuleEngine().Match(processPath, constants.ProcessRuleTypes); err == nil {
			return rule, nil
		}
	}

	ip := metadata.DstIP.String()
	if len(ip) != 0 {
		if rule, err := matcher.GetRuleEngine().Match(ip, constants.IpRuleTypes); err == nil {
			return rule, nil
		}
	}

	host := metadata.Host
	if len(host) != 0 {
		var rule, err = matcher.GetRuleEngine().Match(host, constants.DomainRuleTypes)
		if err == nil {
			return rule, nil
		}
	}

	return nil, fmt.Errorf("no rule found")
}

func NewTun(
	boostDns []string,
	remoteDns []string,
	localDns []string,
	defaultInterfaceName string,
	disableCache bool,
	fakeIpPool *fakeip.Pool,
) (Config, error) {
	if len(boostDns) == 0 || len(remoteDns) == 0 || len(localDns) == 0 {
		return Config{}, fmt.Errorf("dns can't be empty")
	}

	dns.ResetCache()
	dnsConfig, err := NewDnsDistribution(boostDns, remoteDns, localDns, defaultInterfaceName, disableCache, fakeIpPool)
	if err != nil {
		return Config{}, err
	}

	return Config{
		dnsConfig,
	}, nil
}
