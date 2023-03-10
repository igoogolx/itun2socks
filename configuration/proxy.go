package configuration

import (
	"fmt"
	"github.com/Dreamacro/clash/adapter"
	"github.com/gofrs/uuid"
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
		if v["id"] == id {
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

func DeleteProxy(id string) error {
	data, err := Read()
	if err != nil {
		return err
	}
	newProxy := make([]map[string]interface{}, 0)
	for _, v := range data.Proxy {
		if v["id"] != id {
			newProxy = append(newProxy, v)
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

func AddProxies(proxies []map[string]interface{}) ([]map[string]interface{}, error) {
	data, err := Read()
	if err != nil {
		return nil, err
	}
	for _, proxy := range proxies {
		_, err := adapter.ParseProxy(proxy)
		if err != nil {
			return nil, fmt.Errorf("fail to parse proxy,error:%v", err)
		}

		id, err := uuid.NewV4()
		if err != nil {
			return nil, err
		}
		proxy["id"] = id.String()
		data.Proxy = append(data.Proxy, proxy)
	}

	err = Write(data)
	if err != nil {
		return nil, err
	}
	return data.Proxy, nil
}

func AddProxy(proxy map[string]interface{}) (string, error) {
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

func SetClasYamlUrl(url string) error {
	data, err := Read()
	if err != nil {
		return err
	}
	data.ClashYamlUrl = url
	err = Write(data)
	if err != nil {
		return err
	}
	return nil
}

func GetClasYamlUrl() (string, error) {
	data, err := Read()
	if err != nil {
		return "", err
	}
	return data.ClashYamlUrl, nil
}
