package local_server

import (
	"github.com/igoogolx/itun2socks/internal/configuration"
	"strconv"
)

type Cfg struct {
	Addr     string
	AllowLan bool
	Port     int
}

func New(config configuration.LocalServer) Cfg {
	cfg := Cfg{AllowLan: config.AllowLan, Port: config.Port}
	if config.AllowLan {
		cfg.Addr = "0.0.0.0:" + strconv.Itoa(config.Port)
	} else {
		cfg.Addr = "127.0.0.1:" + strconv.Itoa(config.Port)
	}
	return cfg
}
