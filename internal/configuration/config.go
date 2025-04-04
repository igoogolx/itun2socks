package configuration

type Config struct {
	Proxy    []map[string]interface{} `json:"proxy"`
	Selected struct {
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
		DisableCache bool `json:"disableCache"`
		Server       struct {
			Boost  []string `json:"boost"`
			Remote []string `json:"remote"`
			Local  []string `json:"local"`
		} `json:"server"`
		CustomizedOptions []string `json:"customizedOptions"`
	} `json:"dns"`
	Language          string `json:"language,omitempty"`
	BlockQuic         bool   `json:"blockQuic,omitempty"`
	Stack             string `json:"stack"`
	ShouldFindProcess bool   `json:"shouldFindProcess,omitempty"`
	Theme             string `json:"theme,omitempty"`
	AutoConnect       bool   `json:"autoConnect,omitempty"`
	AutoLaunch        bool   `json:"autoLaunch,omitempty"`
	SensitiveInfoMode bool   `json:"sensitiveInfoMode,omitempty"`
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
	AlwaysReset    bool   `json:"alwaysReset,omitempty"`
}
