package configuration_types

type dnsRule struct {
	Address  string   `json:"address"`
	Domains  []string `json:"domains"`
	GeoSites []string `json:"geoSites"`
}

type DnsItem struct {
	Local  dnsRule `json:"local"`
	Remote dnsRule `json:"remote"`
}
