package configuration

type Config struct {
	ClashYamlUrl string                   `json:"clashYamlUrl"`
	Proxy        []map[string]interface{} `json:"proxy"`
	Rule         []RuleCfg                `json:"rule"`
	Selected     struct {
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
	Address  string   `json:"address"`
	Domains  []string `json:"domains"`
	GeoSites []string `json:"geoSites"`
}

type DnsItem struct {
	Local  dnsRule `json:"local"`
	Remote dnsRule `json:"remote"`
}

type IpRule struct {
	Bypass []string `json:"bypass"`
	Proxy  []string `json:"proxy"`
}

type IpItem struct {
	Name         string `json:"name"`
	DefaultProxy bool   `json:"defaultProxy"`
	GeoIps       IpRule `json:"geoIps"`
	GeoSites     IpRule `json:"geoSites"`
	Subnet       IpRule `json:"subnet"`
}

type SettingCfg struct {
	DefaultInterface string      `json:"defaultInterface"`
	TrueProxyServer  string      `json:"trueProxyServer"`
	LocalServer      LocalServer `json:"localServer"`
	Outbound         Outbound    `json:"outbound"`
}

type LocalServer struct {
	Http struct {
		Port    int  `json:"port"`
		Enabled bool `json:"enabled"`
	} `json:"http"`
}

type Outbound struct {
	Mode   string            `json:"mode"`
	Config map[string]string `json:"config"`
}
