package configuration

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/igoogolx/itun2socks/pkg/log"
	"go.uber.org/atomic"
)

var mux sync.RWMutex
var configFilePath = atomic.NewString("")

//go:embed assets/config.json
var defaultConfigContent []byte

func Read() (Config, error) {
	mux.RLock()
	defer mux.RUnlock()
	var err error
	config, err := readFile()
	if err != nil {
		return Config{}, err
	}

	//FIXME: remove in v2
	if len(config.Setting.Stack) == 0 {
		config.Setting.Stack = "gvisor"
	}
	if len(config.Setting.Theme) == 0 {
		config.Setting.Theme = "system"
	}
	if config.Subscriptions == nil {
		config.Subscriptions = []SubscriptionCfg{}
	}

	return *config, nil
}

func Write(c Config) error {
	mux.Lock()
	defer mux.Unlock()
	err := writeFile(c)
	if err != nil {
		return fmt.Errorf("fail to write file: %v", err)
	}
	return nil
}

func Reset() error {
	mux.Lock()
	defer mux.Unlock()
	return write(defaultConfigContent)
}

func SetConfigFilePath(path string) {
	configFilePath.Store(path)
}

func GetConfigFilePath() (string, error) {
	return configFilePath.Load(), nil
}

func readFile() (*Config, error) {
	if !fileExists(configFilePath.Load()) {
		err := write(defaultConfigContent)
		if err != nil {
			return nil, err
		}
		log.Infoln("%s", log.FormatLog(log.ConfigurationPrefix, "created the default config file"))
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
	buf, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		return fmt.Errorf("fail to marchal json, err:%v", err)
	}
	return write(buf)
}

func write(data []byte) error {
	f, err := os.OpenFile(configFilePath.Load(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
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

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
