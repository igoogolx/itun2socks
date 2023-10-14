package constants

type RuleType string

const (
	DistributionBypass RuleType = "bypass"
	DistributionProxy  RuleType = "proxy"
)

type DnsType string

const (
	DistributionLocalDns  DnsType = "local"
	DistributionRemoteDns DnsType = "remote"
	DistributionBoostDns  DnsType = "boost"
)

type RuleConfig string

const (
	RuleIpCidr RuleConfig = "IP-CIDR"
	RuleDomain RuleConfig = "DOMAIN"
)

const DnsPort = "53"

const (
	TunName      = "utun"
	TunLocalAddr = "10.255.0.2"
	TunGateway   = "10.255.0.1"
	TunMask      = "255.255.255.255"
	TunMtu       = 1500
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
