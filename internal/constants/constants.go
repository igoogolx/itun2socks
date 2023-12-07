package constants

type RuleType string

const (
	RuleBypass RuleType = "bypass"
	RuleProxy  RuleType = "proxy"
	RuleReject RuleType = "reject"
)

type DnsType string

const (
	LocalDns  DnsType = "local"
	RemoteDns DnsType = "remote"
	BoostDns  DnsType = "boost"
)

type RuleConfig string

const (
	RuleIpCidr RuleConfig = "IP-CIDR"
	RuleDomain RuleConfig = "DOMAIN"
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
