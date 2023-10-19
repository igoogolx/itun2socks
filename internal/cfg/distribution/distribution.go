package distribution

import (
	"fmt"
	"github.com/igoogolx/itun2socks/internal/cfg/distribution/ruleEngine"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/pkg/log"
	"strings"
)

type Config struct {
	Dns        DnsDistribution
	RuleEngine *ruleEngine.Engine
}

func (c Config) GetDnsServerRule(ip string) (constants.RuleType, error) {
	if strings.Contains(c.Dns.Remote.Address, ip) {
		return constants.RuleProxy, nil
	}
	return constants.RuleBypass, fmt.Errorf("not found")
}

func (c Config) GetDnsRule(ip string) (constants.RuleType, error) {
	result := constants.RuleBypass
	_, cachedRule, ok := GetCachedDnsItem(ip)
	if ok {
		if cachedRule == constants.LocalDns {
			result = constants.RuleBypass
		} else {
			result = constants.RuleProxy
		}
		return result, nil
	}
	return constants.RuleBypass, fmt.Errorf("not found")
}

func (c Config) GetRule(ip string) constants.RuleType {

	result := constants.RuleProxy

	defer func(latestIp string) {
		domain := "unknown"
		dnsRule := "unknown"
		cacheDomain, cachedRule, ok := GetCachedDnsItem(latestIp)
		if ok {
			domain = cacheDomain
			dnsRule = string(cachedRule)
		}
		log.Infoln(log.FormatLog(log.RulePrefix, "ip:%v, rule:%v; domain:%v, rule:%v"), latestIp, result, domain, dnsRule)
	}(ip)

	//dns server
	result, err := c.GetDnsServerRule(ip)
	if err == nil {
		return result
	}

	//dns result
	result, err = c.GetDnsRule(ip)
	if err == nil {
		return result
	}

	rule, err := c.RuleEngine.Match(ip)
	if err == nil {
		return rule.Policy()
	}

	return constants.RuleProxy

}

func (c Config) GetDns(domain string) SubDnsDistribution {
	result := constants.RemoteDns
	if strings.Contains(c.Dns.Remote.Address, domain) {
		result = constants.BoostDns
	} else {
		var rule, err = c.RuleEngine.Match(domain)
		if err == nil {
			if rule.Policy() == constants.RuleBypass {
				result = constants.LocalDns
			} else if rule.Policy() == constants.RuleProxy {
				result = constants.RemoteDns
			}
		}
	}

	switch result {
	case constants.LocalDns:
		return c.Dns.Local
	case constants.RemoteDns:
		return c.Dns.Remote
	default:
		return c.Dns.Boost
	}
}

func New(
	boostDns string,
	remoteDns string,
	localDns string,
	ruleId string,
	rules []string,
	tunInterfaceName string,
	defaultInterfaceName string,
) (Config, error) {
	rEngine, err := ruleEngine.New(ruleId, rules)
	if err != nil {
		return Config{}, err
	}
	dns, err := NewDnsDistribution(boostDns, remoteDns, localDns, tunInterfaceName, defaultInterfaceName)
	if err != nil {
		return Config{}, err
	}
	return Config{
		Dns: dns, RuleEngine: rEngine,
	}, nil
}
