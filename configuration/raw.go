package configuration

import (
	"encoding/json"
	"fmt"
	"github.com/Dreamacro/clash/log"
	configurationTypes "github.com/igoogolx/itun2socks/configuration/configuration-types"
	"reflect"
	"sync"
)

var mux sync.RWMutex
var defaultConfig *configurationTypes.Config

func deepCopy(v interface{}) (interface{}, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	vptr := reflect.New(reflect.TypeOf(v))
	err = json.Unmarshal(data, vptr.Interface())
	if err != nil {
		return nil, err
	}
	return vptr.Elem().Interface(), err
}

func deepCopyConfig(c configurationTypes.Config) (configurationTypes.Config, error) {
	copiedConfig, err := deepCopy(c)
	if err != nil {
		return configurationTypes.Config{}, fmt.Errorf("fail to deep copy config, err:%v", err)
	}
	result, ok := copiedConfig.(configurationTypes.Config)
	if !ok {
		return configurationTypes.Config{}, fmt.Errorf("invald copied config, err:%v", err)
	}

	return result, nil
}

func Read() (configurationTypes.Config, error) {
	mux.RLock()
	defer mux.RUnlock()
	var err error
	if defaultConfig == nil {
		defaultConfig, err = configurationTypes.ReadFile()
		if err != nil {
			return configurationTypes.Config{}, err
		}
	}
	copiedConfig, err := deepCopyConfig(*defaultConfig)
	if err != nil {
		return configurationTypes.Config{}, fmt.Errorf("fail to deep copy config, err:%v", err)
	}
	return copiedConfig, nil
}

func Write(c configurationTypes.Config) error {
	mux.Lock()
	defer mux.Unlock()
	copiedConfig, err := deepCopyConfig(c)
	if err != nil {
		return fmt.Errorf("fail to deep copy config, err:%v", err)
	}
	defaultConfig = &copiedConfig

	go func() {
		err := configurationTypes.WriteFile(copiedConfig)
		if err != nil {
			log.Errorln("fail to write file: %v", err)
		}
	}()
	return nil
}
