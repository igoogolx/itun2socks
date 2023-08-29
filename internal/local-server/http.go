package local_server

import (
	"context"
	"errors"
	"github.com/Dreamacro/clash/log"
	"github.com/elazarl/goproxy"
	"net/http"
)

var server *http.Server

func startHttp(addr string) {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	go func() {
		server = &http.Server{Addr: addr, Handler: proxy}
		err := server.ListenAndServe()
		if err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Errorln("fail to start http local server: %v", err)
			}
		}
	}()
}

func stopHttp() error {
	if server != nil {
		err := server.Shutdown(context.Background())
		return err
	}
	return nil
}
