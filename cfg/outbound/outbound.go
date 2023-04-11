package outbound

import (
	"fmt"
	"github.com/Dreamacro/clash/adapter"
	"github.com/Dreamacro/clash/adapter/outboundgroup"
	"github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/constant/provider"
)

type Option struct {
	Mode    string
	Proxies []map[string]interface{}
	Config  map[string]string
}

func New(option Option) (constant.Proxy, error) {

	var proxy constant.Proxy
	var err error
	var ids = []string{}

	switch option.Mode {
	case "select":
		{
			var selectedProxy map[string]interface{}
			for _, v := range option.Proxies {
				if v["id"] == option.Config["selected"] {
					selectedProxy = v
					break
				}
			}
			if selectedProxy == nil {
				return nil, fmt.Errorf("error getting seleted proxyConfig, id:%v", option.Config["selected"])
			}
			proxy, err = adapter.ParseProxy(selectedProxy)
			if err != nil {
				return nil, err
			}
			break
		}
	case "auto":
		{
			proxyMap := map[string]constant.Proxy{}
			for _, v := range option.Proxies {
				p, err := adapter.ParseProxy(v)
				if err != nil {
					return nil, fmt.Errorf("fail to parse proxy: %v", err)
				}
				proxyMap[v["id"].(string)] = p
				ids = append(ids, v["id"].(string))
			}
			proxyGroupConfig := map[string]any{
				"name":     "auto",
				"proxies":  ids,
				"interval": 300,
				"url":      option.Config["url"],
				"type":     option.Config["type"],
			}

			proxyGroup, err := outboundgroup.ParseProxyGroup(proxyGroupConfig, proxyMap, map[string]provider.ProxyProvider{})
			if err != nil {
				return nil, fmt.Errorf("fail to parse proxy group: %v", err)
			}
			proxy = adapter.NewProxy(proxyGroup)
			break
		}

	default:
		return nil, fmt.Errorf("unsupported outbound Mode: %v", option.Mode)

	}

	return proxy, nil

}
