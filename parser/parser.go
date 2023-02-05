package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Dreamacro/clash/adapter"
	"github.com/Dreamacro/clash/adapter/outbound"
	"github.com/Dreamacro/clash/constant"
)

func ParseOption(mapping map[string]interface{}) (interface{}, error) {
	var option interface{}
	var err error
	mapping["port"] = int(mapping["port"].(float64))
	proxyType, existType := mapping["type"].(string)
	if !existType {
		return nil, errors.New("missing type")
	}
	rawMapping, err := json.Marshal(mapping)
	if err != nil {
		return nil, errors.New("invalid mapping")
	}
	switch proxyType {
	case "ss":
		ssOption := &ShadowSocksOption{}
		err = json.Unmarshal(rawMapping, ssOption)
		if err != nil {
			break
		}
		option = ssOption
	case "socks5":
		socksOption := &Socks5Option{}
		err = json.Unmarshal(rawMapping, socksOption)
		if err != nil {
			break
		}
		option = socksOption
	default:
		return nil, fmt.Errorf("unsupport proxy type: %s", proxyType)
	}

	if err != nil {
		return nil, err
	}

	return option, nil
}

func convert(mapping map[string]interface{}) (interface{}, error) {
	proxyType, existType := mapping["type"].(string)
	if !existType {
		return nil, fmt.Errorf("missing type")
	}
	option, err := ParseOption(mapping)
	if err != nil {
		return nil, err
	}
	var proxy interface{}
	switch proxyType {
	case "ss":
		ssOption := option.(*ShadowSocksOption)
		proxy = outbound.ShadowSocksOption{
			Name:       ssOption.Name,
			Server:     ssOption.Server,
			Port:       ssOption.Port,
			Password:   ssOption.Password,
			Cipher:     ssOption.Cipher,
			UDP:        ssOption.UDP,
			Plugin:     ssOption.Plugin,
			PluginOpts: ssOption.PluginOpts,
		}
	case "socks5":
		socksOption := option.(*Socks5Option)
		proxy = outbound.Socks5Option{
			Name:     socksOption.Name,
			Server:   socksOption.Server,
			Port:     socksOption.Port,
			UserName: socksOption.UserName,
			Password: socksOption.Password,
			UDP:      socksOption.UDP,
		}
	default:
		return nil, fmt.Errorf("unsupport proxy type: %s", proxyType)
	}
	return proxy, nil
}

func ParseProxy(mapping map[string]interface{}) (*adapter.Proxy, error) {
	proxyType, existType := mapping["type"].(string)
	if !existType {
		return nil, fmt.Errorf("missing type")
	}
	option, err := convert(mapping)
	if err != nil {
		return nil, err
	}
	var (
		proxy constant.ProxyAdapter
	)
	switch proxyType {
	case "ss":
		ssOption := option.(outbound.ShadowSocksOption)
		proxy, err = outbound.NewShadowSocks(ssOption)
	case "socks5":
		socksOption := option.(outbound.Socks5Option)
		proxy = outbound.NewSocks5(socksOption)
	default:
		return nil, fmt.Errorf("unsupport proxy type: %s", proxyType)
	}

	if err != nil {
		return nil, err
	}

	return adapter.NewProxy(proxy), nil
}
