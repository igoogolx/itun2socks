package configuration_types

type IPRule struct {
	Bypass []string `json:"bypass"`
	Proxy  []string `json:"proxy"`
}

type IpItem struct {
	Name         string `json:"name"`
	DefaultProxy bool   `json:"defaultProxy"`
	GeoIps       IPRule `json:"geoIps"`
	GeoSites     IPRule `json:"geoSites"`
	Subnet       IPRule `json:"subnet"`
}
