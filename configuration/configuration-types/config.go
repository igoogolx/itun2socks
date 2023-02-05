package configuration_types

import (
	"encoding/json"
	"fmt"
	"github.com/Dreamacro/clash/log"
	"go.uber.org/atomic"
	"os"
	"sync"
)

type Config struct {
	Proxy    []map[string]interface{} `json:"proxy"`
	Rule     []RuleCfg                `json:"rule"`
	Selected struct {
		Proxy string `json:"proxy"`
		Rule  string `json:"rule"`
	} `json:"selected"`
	Setting SettingCfg `json:"setting"`
}

type RuleCfg struct {
	Id   string  `json:"id"`
	Name string  `json:"name"`
	Dns  DnsItem `json:"dns"`
	Ip   IpItem  `json:"ip"`
}

var mutex sync.RWMutex
var ConfigFilePath = atomic.NewString("")

func ReadFile() (*Config, error) {
	mutex.RLock()
	defer mutex.RUnlock()
	c := &Config{}
	data, err := os.ReadFile(ConfigFilePath.Load())
	if err != nil {
		return nil, fmt.Errorf("fail to read config file, path:%v, err:%v", ConfigFilePath.Load(), err)
	}
	err = json.Unmarshal(data, c)
	if err != nil {
		return nil, fmt.Errorf("fail to parse config file, path:%v, err:%v", ConfigFilePath.Load(), err)
	}
	return c, nil
}

func WriteFile(config Config) error {
	mutex.Lock()
	defer mutex.Unlock()
	f, err := os.OpenFile(ConfigFilePath.Load(), os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Warnln("fail to close file: %v, err: %v", ConfigFilePath.Load(), err)
		}
	}(f)
	if err != nil {
		return fmt.Errorf("fail to open file:%v, err:%v", ConfigFilePath.Load(), err)
	}
	buf, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		return fmt.Errorf("fail to marchal json, err:%v", err)
	}
	_, err = f.Write(buf)
	if err != nil {
		return fmt.Errorf("fail to write file:%v, err:%v", ConfigFilePath.Load(), err)
	}
	return nil
}
