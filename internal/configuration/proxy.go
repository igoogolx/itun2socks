package configuration

import (
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/Dreamacro/clash/adapter"
	"github.com/gofrs/uuid/v5"
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

func GetSubscriptions() ([]SubscriptionCfg, error) {
	data, err := Read()
	if err != nil {
		return nil, err
	}
	return data.Subscriptions, nil
}

func DeleteSubscription(id string) error {
	data, err := Read()
	if err != nil {
		return err
	}
	var newSubscriptions = make([]SubscriptionCfg, 0)
	for _, s := range data.Subscriptions {
		if s.Id != id {
			newSubscriptions = append(newSubscriptions, s)
		}
	}

	newProxy := make([]map[string]interface{}, 0)

	for _, v := range data.Proxy {
		subscriptionId, ok := v["subscription"].(string)
		if !ok {
			newProxy = append(newProxy, v)
			continue
		}
		if subscriptionId != id {
			newProxy = append(newProxy, v)
		}
	}

	data.Subscriptions = newSubscriptions
	data.Proxy = newProxy

	err = Write(data)
	if err != nil {
		return err
	}

	return nil
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

func AddSubscription(proxies []map[string]interface{}, subscriptionUrl string, subscriptionName string, subscriptionRemark string) ([]map[string]interface{}, []SubscriptionCfg, error) {
	data, err := Read()
	if err != nil {
		return nil, nil, err
	}

	subscriptionUuid, err := uuid.NewV4()
	if err != nil {
		return nil, nil, err
	}

	subscriptionId := subscriptionUuid.String()

	newProxyWithIds := make([]map[string]interface{}, 0)
	for _, proxy := range proxies {
		_, err := adapter.ParseProxy(proxy)
		if err != nil {
			return nil, nil, fmt.Errorf("fail to parse proxy,error:%v", err)
		}

		if _, ok := checkIsValidStr(proxy["id"]); !ok {
			id, err := uuid.NewV4()
			if err != nil {
				return nil, nil, err
			}
			proxy["id"] = id.String()
		}

		proxy["subscription"] = subscriptionId
		newProxyWithIds = append(newProxyWithIds, proxy)
	}

	newSubscription := SubscriptionCfg{Id: subscriptionId, Name: subscriptionName, Remark: subscriptionRemark, Url: subscriptionUrl}
	data.Proxy = append(data.Proxy, newProxyWithIds...)
	data.Subscriptions = append(data.Subscriptions, newSubscription)

	err = Write(data)
	if err != nil {
		return nil, nil, err
	}
	return data.Proxy, data.Subscriptions, nil
}

func UpdateSubscription(subscription SubscriptionCfg) error {
	c, err := Read()
	if err != nil {
		return err
	}
	for i, v := range c.Subscriptions {
		if v.Id == subscription.Id {
			c.Subscriptions[i] = subscription
			break
		}
	}
	return Write(c)
}

func UpdateSubscriptionProxies(subscriptionId string, proxies []map[string]interface{}) ([]map[string]interface{}, error) {
	c, err := Read()
	if err != nil {
		return nil, err
	}

	newProxies := make([]map[string]interface{}, 0)

	for _, p := range c.Proxy {
		if value, ok := p["subscription"].(string); ok && value == subscriptionId {
			continue
		}
		newProxies = append(newProxies, p)
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
		proxy["subscription"] = subscriptionId

		newProxies = append(newProxies, proxy)
	}

	c.Proxy = newProxies

	err = Write(c)
	if err != nil {
		return nil, err
	}

	return newProxies, nil
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
