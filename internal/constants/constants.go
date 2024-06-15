package constants

import "runtime"

type Policy string

const (
	PolicyDirect Policy = "DIRECT"
	PolicyProxy  Policy = "PROXY"
	PolicyReject Policy = "REJECT"
)

type DnsType string

const (
	LocalDns  DnsType = "local"
	RemoteDns DnsType = "remote"
	BoostDns  DnsType = "boost"
)

type RuleType string

const (
	RuleIpCidr        RuleType = "IP-CIDR"
	RuleDomain        RuleType = "DOMAIN"
	RuleDomainKeyword RuleType = "DOMAIN-KEYWORD"
	RuleDomainRegex   RuleType = "DOMAIN-REGEX"
	RuleDomainSuffix  RuleType = "DOMAIN-SUFFIX"
	RuleProcess       RuleType = "PROCESS"
	RuleBuiltIn       RuleType = "BUILT-IN"
)

var (
	IpRuleTypes      = []RuleType{RuleIpCidr}
	DomainRuleTypes  = []RuleType{RuleDomain, RuleDomainSuffix, RuleDomainRegex}
	ProcessRuleTypes = []RuleType{RuleProcess}
)

const DnsPort = "53"

func TunName() string {
	if runtime.GOOS == "windows" {
		//It is tricky. Utun is "bigger" than Wi-Fi so that ethernet icon is displayed on the task bar.
		return "ztun"
	}
	return "utun"
}

const (
	TunGateway = "10.255.0.1/30"
	TunMtu     = 1500
)

const (
	DbFileName = "config.json"
)

const (
	DefaultHubPort = 9000
)

var (
	Version   = "undefined"
	BuildTime = "undefined"
)
