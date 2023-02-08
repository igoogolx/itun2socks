package configuration

import (
	"encoding/json"
	"fmt"
	"github.com/Dreamacro/clash/log"
	"go.uber.org/atomic"
	"os"
	"reflect"
	"sync"
)

var mux sync.RWMutex
var defaultConfig *Config

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

func deepCopyConfig(c Config) (Config, error) {
	copiedConfig, err := deepCopy(c)
	if err != nil {
		return Config{}, fmt.Errorf("fail to deep copy config, err:%v", err)
	}
	result, ok := copiedConfig.(Config)
	if !ok {
		return Config{}, fmt.Errorf("invald copied config, err:%v", err)
	}

	return result, nil
}

func Read() (Config, error) {
	mux.RLock()
	defer mux.RUnlock()
	var err error
	if defaultConfig == nil {
		defaultConfig, err = readFile()
		if err != nil {
			return Config{}, err
		}
	}
	copiedConfig, err := deepCopyConfig(*defaultConfig)
	if err != nil {
		return Config{}, fmt.Errorf("fail to deep copy config, err:%v", err)
	}
	return copiedConfig, nil
}

func Write(c Config) error {
	mux.Lock()
	defer mux.Unlock()
	copiedConfig, err := deepCopyConfig(c)
	if err != nil {
		return fmt.Errorf("fail to deep copy config, err:%v", err)
	}
	defaultConfig = &copiedConfig

	go func() {
		err := writeFile(copiedConfig)
		if err != nil {
			log.Errorln("fail to write file: %v", err)
		}
	}()
	return nil
}

var fileMutex sync.RWMutex
var ConfigFilePath = atomic.NewString("")

func readFile() (*Config, error) {
	fileMutex.RLock()
	defer fileMutex.RUnlock()
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

func writeFile(config Config) error {
	fileMutex.Lock()
	defer fileMutex.Unlock()
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
