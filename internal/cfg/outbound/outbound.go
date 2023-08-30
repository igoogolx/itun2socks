package outbound

import (
	"fmt"
	"github.com/Dreamacro/clash/adapter"
	"github.com/Dreamacro/clash/adapter/outboundgroup"
	"github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/constant/provider"
	"github.com/igoogolx/itun2socks/internal/configuration"
)

type Option struct {
	AutoMode      configuration.AutoMode
	Proxies       []map[string]interface{}
	SelectedProxy string
}

func New(option Option) (constant.Proxy, error) {

	var proxy constant.Proxy
	var err error
	var ids []string

	if option.AutoMode.Enabled {

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
				"url":      option.AutoMode.Url,
				"type":     option.AutoMode.Type,
			}

			proxyGroup, err := outboundgroup.ParseProxyGroup(proxyGroupConfig, proxyMap, map[string]provider.ProxyProvider{})
			if err != nil {
				return nil, fmt.Errorf("fail to parse proxy group: %v", err)
			}
			proxy = adapter.NewProxy(proxyGroup)
		}
	} else {
		var selectedProxy map[string]interface{}
		for _, v := range option.Proxies {
			if v["id"] == option.SelectedProxy {
				selectedProxy = v
				break
			}
		}
		if selectedProxy == nil {
			return nil, fmt.Errorf("error getting seleted proxyConfig, id:%v", option.SelectedProxy)
		}
		proxy, err = adapter.ParseProxy(selectedProxy)
		if err != nil {
			return nil, err
		}
	}

	return proxy, nil

}
