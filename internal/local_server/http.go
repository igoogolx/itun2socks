package local_server

import (
	"context"
	"errors"
	"github.com/elazarl/goproxy"
	"github.com/igoogolx/itun2socks/pkg/log"
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
				log.Errorln(log.FormatLog(log.HubPrefix, "fail to start http local server: %v"), err)
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
