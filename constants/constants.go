package constants

import "path/filepath"

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
)

const CacheSize = 1000

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
	HomeDir   = ""
	Version   = "undefined"
	BuildTime = "undefined"
)

var LogFile = filepath.Join("logs", "lux.log")
