package parser

type BaseInfo struct {
	Id   string `json:"id,omitempty"`
	Type string `json:"type"`
}

type ShadowSocksOption struct {
	BaseInfo
	Name       string                 `json:"name"`
	Server     string                 `json:"server"`
	Port       int                    `json:"port"`
	Password   string                 `json:"password"`
	Cipher     string                 `json:"method"`
	UDP        bool                   `json:"udp,omitempty"`
	Plugin     string                 `json:"plugin,omitempty"`
	PluginOpts map[string]interface{} `json:"pluginOpts,omitempty"`
}

type Socks5Option struct {
	BaseInfo
	Name           string `json:"name"`
	Server         string `json:"server"`
	Port           int    `json:"port"`
	UserName       string `json:"username,omitempty"`
	Password       string `json:"password,omitempty"`
	TLS            bool   `json:"tls,omitempty"`
	UDP            bool   `json:"udp,omitempty"`
	SkipCertVerify bool   `json:"skipCertVerify,omitempty"`
}
