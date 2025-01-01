package configuration

import (
	"fmt"
	"github.com/Dreamacro/clash/adapter"
	"github.com/gofrs/uuid/v5"
	"slices"
	"strings"
	"sync"
)

func GetSelectedProxy() (map[string]interface{}, error) {
	data, err := Read()
	if err != nil {
		return nil, err
	}
	return GetProxy(data.Selected.Proxy)
}

func GetProxy(id string) (map[string]interface{}, error) {
	data, err := Read()
	if err != nil {
		return nil, err
	}
	for _, v := range data.Proxy {
		if v["id"] == id && v != nil {
			return v, nil
		}
	}
	return nil, fmt.Errorf("error getting selected proxy,id:%v,err:%v", id, err)
}

func GetProxies() ([]map[string]interface{}, error) {
	data, err := Read()
	if err != nil {
		return nil, err
	}
	return data.Proxy, nil
}

func DeleteProxies(ids []string) error {
	data, err := Read()
	if err != nil {
		return err
	}
	newProxy := make([]map[string]interface{}, 0)
	for _, v := range data.Proxy {
		if id, ok := v["id"].(string); ok && len(id) != 0 {
			if !slices.Contains(ids, id) {
				newProxy = append(newProxy, v)
			}
		}
	}
	data.Proxy = newProxy
	err = Write(data)
	if err != nil {
		return err
	}
	return nil
}

func DeleteAllProxies() error {
	data, err := Read()
	if err != nil {
		return err
	}
	data.Proxy = make([]map[string]interface{}, 0)
	err = Write(data)
	return err
}

func UpdateProxy(id string, proxy map[string]interface{}) error {
	_, err := adapter.ParseProxy(proxy)
	if err != nil {
		return fmt.Errorf("fail to update proxy,error:%v", err)
	}
	c, err := Read()
	if err != nil {
		return err
	}
	for i, v := range c.Proxy {
		if v["id"] == id {
			c.Proxy[i] = proxy
			break
		}
	}
	err = Write(c)
	if err != nil {
		return err
	}
	return nil
}

func checkIsValidStr(value interface{}) (string, bool) {
	str, ok := value.(string)
	if !ok {
		return "", false
	}
	if len(strings.TrimSpace(str)) == 0 {
		return "", false
	}
	return str, true
}

func AddProxies(proxies []map[string]interface{}, subscriptionUrl string) ([]map[string]interface{}, error) {
	data, err := Read()
	if err != nil {
		return nil, err
	}

	newProxyWithIds := make([]map[string]interface{}, 0)
	idNs := uuid.NewV5(uuid.Nil, subscriptionUrl)
	for _, proxy := range proxies {
		_, err := adapter.ParseProxy(proxy)
		if err != nil {
			return nil, fmt.Errorf("fail to parse proxy,error:%v", err)
		}

		if proxyName, ok := checkIsValidStr(proxy["name"]); ok {
			id := uuid.NewV5(idNs, proxyName)
			proxy["id"] = id.String()
		}

		if _, ok := checkIsValidStr(proxy["id"]); !ok {
			id, err := uuid.NewV4()
			if err != nil {
				return nil, err
			}
			proxy["id"] = id.String()
		}

		proxy["subscriptionUrl"] = subscriptionUrl
		newProxyWithIds = append(newProxyWithIds, proxy)
	}

	targetIndex := slices.IndexFunc(data.Proxy, func(p map[string]interface{}) bool {
		return p["subscriptionUrl"] == subscriptionUrl
	})

	newProxy := make([]map[string]interface{}, 0)

	for _, v := range data.Proxy {
		if v["subscriptionUrl"] != subscriptionUrl {
			newProxy = append(newProxy, v)
		}
	}

	if targetIndex != -1 && targetIndex < len(newProxy) {
		newProxy = slices.Insert(newProxy, targetIndex, newProxyWithIds...)
	} else {
		newProxy = append(newProxy, newProxyWithIds...)
	}
	data.Proxy = newProxy
	err = Write(data)
	if err != nil {
		return nil, err
	}
	return data.Proxy, nil
}

var addMux sync.Mutex

func AddProxy(proxy map[string]interface{}) (string, error) {
	addMux.Lock()
	defer addMux.Unlock()
	_, err := adapter.ParseProxy(proxy)
	if err != nil {
		return "", fmt.Errorf("fail to parse proxy,error:%v", err)
	}
	data, err := Read()
	if err != nil {
		return "", err
	}
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	proxy["id"] = id.String()
	data.Proxy = append(data.Proxy, proxy)
	err = Write(data)
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
