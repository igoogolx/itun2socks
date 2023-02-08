package configuration_types

type Config struct {
	Proxy    []map[string]interface{} `json:"proxy"`
	Rule     []RuleCfg                `json:"rule"`
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
