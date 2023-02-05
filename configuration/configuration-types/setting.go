package configuration_types

type SettingCfg struct {
	DefaultInterface string      `json:"defaultInterface"`
	TrueProxyServer  string      `json:"trueProxyServer"`
	LocalServer      LocalServer `json:"localServer"`
}

type LocalServer struct {
	Http struct {
		Port    int  `json:"port"`
		Enabled bool `json:"enabled"`
	} `json:"http"`
}
