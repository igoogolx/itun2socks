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

type RuleConfig string

const (
	RuleIpCidr        RuleConfig = "IP-CIDR"
	RuleDomain        RuleConfig = "DOMAIN"
	RuleDomainKeyword RuleConfig = "DOMAIN-KEYWORD"
	RuleDomainRegex   RuleConfig = "DOMAIN-REGEX"
	RuleDomainSuffix  RuleConfig = "DOMAIN-SUFFIX"
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
