package outbound

import (
	"fmt"
	"github.com/metacubex/mihomo/adapter"
	"github.com/metacubex/mihomo/adapter/outboundgroup"
	"github.com/metacubex/mihomo/constant/provider"

	"github.com/igoogolx/itun2socks/internal/configuration"
	metaC "github.com/metacubex/mihomo/constant"
)

type Option struct {
	AutoMode      configuration.AutoMode
	Proxies       []map[string]any
	SelectedProxy string
}

func New(option Option) (metaC.Proxy, error) {

	var proxy metaC.Proxy
	var err error
	var ids []string
	var allProxies []string

	if option.AutoMode.Enabled {

		{
			proxyMap := map[string]metaC.Proxy{}
			for _, v := range option.Proxies {
				p, err := adapter.ParseProxy(v)
				if err != nil {
					return nil, fmt.Errorf("fail to parse proxy: %v", err)
				}
				proxyMap[v["id"].(string)] = p
				ids = append(ids, v["id"].(string))
				allProxies = append(allProxies, v["name"].(string))
			}
			proxyGroupConfig := map[string]any{
				"name":     "auto",
				"proxies":  ids,
				"interval": 300,
				"url":      option.AutoMode.Url,
				"type":     option.AutoMode.Type,
			}

			proxyGroup, err := outboundgroup.ParseProxyGroup(proxyGroupConfig, proxyMap, map[string]provider.ProxyProvider{}, allProxies, []string{})
			if err != nil {
				return nil, fmt.Errorf("fail to parse proxy group: %v", err)
			}
			proxy = adapter.NewProxy(proxyGroup)
		}
	} else {
		var selectedProxy map[string]any
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
