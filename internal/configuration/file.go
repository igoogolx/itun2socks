package configuration

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/igoogolx/itun2socks/pkg/log"
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
	err = writeFile(copiedConfig)
	if err != nil {
		return fmt.Errorf("fail to write file: %v", err)
	}
	return nil
}

var fileMutex sync.Mutex
var configFilePath = atomic.NewString("")

func SetConfigFilePath(path string) {
	configFilePath.Store(path)
}

func Reset() error {
	fileMutex.Lock()
	defer fileMutex.Unlock()
	return write(defaultConfigContent)
}

func GetConfigFilePath() (string, error) {
	return configFilePath.Load(), nil
}

//go:embed assets/config.json
var defaultConfigContent []byte

func readFile() (*Config, error) {
	fileMutex.Lock()
	defer fileMutex.Unlock()
	if !fileExists(configFilePath.Load()) {
		err := write(defaultConfigContent)
		if err != nil {
			return nil, err
		}
		log.Infoln(log.FormatLog(log.ConfigurationPrefix, "created the default config file"))
	}
	c := &Config{}
	data, err := os.ReadFile(configFilePath.Load())
	if err != nil {
		return nil, fmt.Errorf("fail to read config file, path:%v, err:%v", configFilePath.Load(), err)
	}
	err = json.Unmarshal(data, c)
	if err != nil {
		log.Errorln("fail to parse config file, path:%v, err:%v, using default file", configFilePath.Load(), err)
		err = json.Unmarshal(defaultConfigContent, c)
		if err != nil {
			return nil, fmt.Errorf("fail to parse defalut config file: %v", err)
		}
	}
	return c, nil
}

func writeFile(config Config) error {
	fileMutex.Lock()
	defer fileMutex.Unlock()
	buf, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		return fmt.Errorf("fail to marchal json, err:%v", err)
	}
	return write(buf)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func write(data []byte) error {
	f, err := os.Create(configFilePath.Load())
	if err != nil {
		return fmt.Errorf("fail to open file:%v, err:%v", configFilePath.Load(), err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Warnln(log.FormatLog(log.ConfigurationPrefix, "fail to close file: %v, err: %v"), configFilePath.Load(), err)
		}
	}(f)
	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("fail to write file:%v, err:%v", configFilePath.Load(), err)
	}
	return nil
}
