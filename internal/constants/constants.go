package constants

type Policy string

const (
	PolicyDirect Policy = "DIRECT"
	PolicyProxy  Policy = "PROXY"
	PolicyReject Policy = "REJECT"
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
	DomainRuleTypes  = []RuleType{RuleDomain, RuleDomainSuffix, RuleDomainRegex, RuleDomainKeyword}
	ProcessRuleTypes = []RuleType{RuleProcess}
)

const DnsPort = "53"

func TunName() string {
	return "utun"
}

const (
	TunGateway  = "10.255.0.1/30"
	TunMtu      = 1500
	HijackedDns = "10.255.0.2"
)

const (
	DbFileName = "config.json"
)

var (
	Version   = "undefined"
	BuildTime = "undefined"
)

var FakeIpRange = "198.18.0.1/16"
