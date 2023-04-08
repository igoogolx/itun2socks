package outbound

import (
	"fmt"
	"github.com/Dreamacro/clash/adapter"
	"github.com/Dreamacro/clash/adapter/outboundgroup"
	"github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/constant/provider"
)

func New(proxyConfig []map[string]interface{}, selected string, fallback bool) (constant.Proxy, error) {

	var proxy constant.Proxy
	var err error
	var ids = []string{}

	if fallback {
		var proxyMap = map[string]constant.Proxy{}
		for _, v := range proxyConfig {
			p, err := adapter.ParseProxy(v)
			if err != nil {
				return nil, fmt.Errorf("fail to parse proxy: %v", err)
			}
			proxyMap[v["id"].(string)] = p
			ids = append(ids, v["id"].(string))
		}
		proxyGroupConfig := map[string]any{
			"name":     "fallback-auto",
			"type":     "url-test",
			"proxies":  ids,
			"url":      "https://www.google.com/",
			"interval": 300,
		}

		proxyGroup, err := outboundgroup.ParseProxyGroup(proxyGroupConfig, proxyMap, map[string]provider.ProxyProvider{})
		if err != nil {
			return nil, fmt.Errorf("fail to parse proxy group: %v", err)
		}
		proxy = adapter.NewProxy(proxyGroup)
	} else {
		var selectedProxy map[string]interface{}
		for _, v := range proxyConfig {
			if v["id"] == selected {
				selectedProxy = v
				break
			}
		}
		if selectedProxy == nil {
			return nil, fmt.Errorf("error getting seleted proxyConfig, id:%v", selected)
		}
		proxy, err = adapter.ParseProxy(selectedProxy)
		if err != nil {
			return nil, err
		}
	}

	return proxy, nil

}
