package configuration

type Config struct {
	ClashYamlUrl string                   `json:"clashYamlUrl"`
	Proxy        []map[string]interface{} `json:"proxy"`
	Selected     struct {
		Proxy string `json:"proxy"`
		Rule  string `json:"rule"`
	} `json:"selected"`
	Setting SettingCfg `json:"setting"`
	Rules   []string   `json:"rules"`
}

type SettingCfg struct {
	Mode             string `json:"mode"`
	DefaultInterface string `json:"defaultInterface"`
	LocalServer      `json:"localServer"`
	AutoMode         `json:"autoMode"`
	HijackDns        `json:"hijackDns"`
	Dns              struct {
		Boost  DnsServer `json:"boost"`
		Remote DnsServer `json:"remote"`
		Local  DnsServer `json:"local"`
	} `json:"dns"`
	Language  string `json:"language,omitempty"`
	BlockQuic bool   `json:"blockQuic,omitempty"`
}

type DnsServer struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type LocalServer struct {
	Port     int  `json:"port"`
	AllowLan bool `json:"allowLan"`
}

type AutoMode struct {
	Enabled bool   `json:"enabled"`
	Type    string `json:"type"`
	Url     string `json:"url"`
}

type HijackDns struct {
	Enabled        bool   `json:"enabled"`
	NetworkService string `json:"networkService"`
}
