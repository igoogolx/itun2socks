package configuration

type Config struct {
	ClashYamlUrl string                   `json:"clashYamlUrl"`
	Proxy        []map[string]interface{} `json:"proxy"`
	Rules        []RuleCfg                `json:"rules"`

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
		Boost  DnsServer `json:"boost"`
		Remote DnsServer `json:"remote"`
		Local  DnsServer `json:"local"`
	} `json:"dns"`
}

type DnsServer struct {
	Type  string `json:"type"`
	Value string `json:"value"`
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
