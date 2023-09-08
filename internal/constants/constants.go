package constants

type IpRule int

const (
	DistributionBypass   IpRule = 0
	DistributionProxy    IpRule = 1
	DistributionNotFound IpRule = 2
)

type DnsRule string

const (
	DistributionLocalDns  DnsRule = "local"
	DistributionRemoteDns DnsRule = "remote"
	DistributionBoostDns  DnsRule = "boost"
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
