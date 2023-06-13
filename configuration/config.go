package configuration

type Config struct {
	ClashYamlUrl string                   `json:"clashYamlUrl"`
	Proxy        []map[string]interface{} `json:"proxy"`
	Rule         []RuleCfg                `json:"rule"`

	Selected struct {
		Proxy string `json:"proxy"`
		Rule  string `json:"rule"`
	} `json:"selected"`
	Setting SettingCfg `json:"setting"`
}

type RuleCfg struct {
	Id   string  `json:"id"`
	Name string  `json:"name"`
	Dns  DnsItem `json:"dns"`
	Ip   IpItem  `json:"ip"`
}

type dnsRule struct {
	Local  []string `json:"local"`
	Remote []string `json:"remote"`
}

type DnsItem struct {
	Domains  dnsRule `json:"domains"`
	GeoSites dnsRule `json:"geoSites"`
}

type IpRule struct {
	Bypass []string `json:"bypass"`
	Proxy  []string `json:"proxy"`
}

type IpItem struct {
	Name         string `json:"name"`
	DefaultProxy bool   `json:"defaultProxy"`
	GeoIps       IpRule `json:"geoIps"`
	Subnet       IpRule `json:"subnet"`
}

type SettingCfg struct {
	DefaultInterface string `json:"defaultInterface"`
	TrueProxyServer  string `json:"trueProxyServer"`
	LocalServer      `json:"localServer"`
	AutoMode         `json:"autoMode"`
	Dns              struct {
		Boost  string `json:"boost"`
		Remote string `json:"remote"`
		Local  string `json:"local"`
	} `json:"dns"`
}

type LocalServer struct {
	Http struct {
		Port    int  `json:"port"`
		Enabled bool `json:"enabled"`
	} `json:"http"`
}

type AutoMode struct {
	Enabled bool   `json:"enabled"`
	Type    string `json:"type"`
	Url     string `json:"url"`
}
