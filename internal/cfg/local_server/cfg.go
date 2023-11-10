package local_server

import (
	"github.com/igoogolx/itun2socks/internal/configuration"
	"strconv"
)

type Cfg struct {
	HttpAddr string
}

func New(config configuration.LocalServer) Cfg {
	cfg := Cfg{}
	if config.Http.Enabled {
		cfg.HttpAddr = "127.0.0.1:" + strconv.Itoa(config.Http.Port)
	} else {
		cfg.HttpAddr = "127.0.0.1:" + strconv.Itoa(config.Http.Port)
	}
	return cfg
}
