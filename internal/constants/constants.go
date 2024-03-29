package constants

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
)

const DnsPort = "53"

const (
	TunName    = "utun"
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
