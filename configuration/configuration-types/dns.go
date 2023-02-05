package configuration_types

type DnsRule struct {
	Address  string   `json:"address"`
	Domains  []string `json:"domains"`
	GeoSites []string `json:"geoSites"`
}

type DnsItem struct {
	Local  DnsRule `json:"local"`
	Remote DnsRule `json:"remote"`
}
