package routes

import (
	"context"
	"github.com/Dreamacro/clash/adapter"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	db2 "github.com/igoogolx/itun2socks/configuration"
	"github.com/igoogolx/itun2socks/tunnel"
	"net/http"
	"sync"
	"time"
)

var (
	defaultDelayTimeout = 5 * time.Second
	defaultDelayTestUrl = "https://www.google.com"
)

func proxyRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getProxies)
	r.Put("/", addProxy)
	r.Delete("/{proxyId}", deleteProxy)
	r.Post("/{proxyId}", updateProxy)
	r.Get("/delay/{proxyId}", getProxyDelay)
	r.Get("/delays", getProxiesDelay)
	r.Get("/udp-test/{proxyId}", testProxyUdp)
	return r
}

func testProxyUdp(w http.ResponseWriter, r *http.Request) {
	proxyId := chi.URLParam(r, "proxyId")
	if proxyId == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	url := chi.URLParam(r, "url")
	if url == "" {
		url = defaultDelayTestUrl
	}
	proxyOption, err := db2.GetProxy(proxyId)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	p, err := adapter.ParseProxy(proxyOption)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	metadata, err := tunnel.CreateMetadata("", "8.8.8.8:53", C.UDP)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	pc, _ := p.ListenPacketContext(context.Background(), metadata)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	res := UdpTest(pc, "8.8.8.8:53")

	render.JSON(w, r, render.M{
		"result": res,
	})
}

func getProxyDelay(w http.ResponseWriter, r *http.Request) {
	proxyId := chi.URLParam(r, "proxyId")
	if proxyId == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	url := chi.URLParam(r, "url")
	if url == "" {
		url = defaultDelayTestUrl
	}
	proxyOption, err := db2.GetProxy(proxyId)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	p, err := adapter.ParseProxy(proxyOption)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultDelayTimeout)
	defer cancel()
	delay, err := p.URLTest(ctx, url)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	render.JSON(w, r, render.M{
		"delay": delay,
	})
}

func getProxiesDelay(w http.ResponseWriter, r *http.Request) {
	url := chi.URLParam(r, "url")
	if url == "" {
		url = defaultDelayTestUrl
	}
	proxyConfigs, err := db2.GetProxies()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}

	proxies := make([]C.Proxy, len(proxyConfigs))
	for _, proxy := range proxyConfigs {
		p, err := adapter.ParseProxy(proxy)
		if err != nil {
			continue
		}
		proxies = append(proxies, p)
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultDelayTimeout)
	defer cancel()

	wg := &sync.WaitGroup{}

	type Delay struct {
		Id    string `json:"id"`
		Value uint16 `json:"value"`
	}
	delays := make([]Delay, 0)
	var m sync.Mutex
	for i, proxy := range proxies {
		wg.Add(1)
		i := i
		go func(p C.Proxy) {
			delay, err := p.URLTest(ctx, url)
			if err != nil {
				log.Errorln("error:%v", err)
			}
			m.Lock()
			currentProxy := proxyConfigs[i]
			delays = append(delays, Delay{Id: currentProxy["id"].(string), Value: delay})
			m.Unlock()
			wg.Done()
		}(proxy)
	}
	wg.Wait()
	render.JSON(w, r, render.M{
		"delays": delays,
	})
}

func getProxies(w http.ResponseWriter, r *http.Request) {
	proxiesMap, err := db2.GetProxies()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	selectedId, err := db2.GetSelectedId("proxy")
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	proxies := make([]interface{}, 0)
	for _, proxy := range proxiesMap {
		proxies = append(proxies, proxy)
	}
	render.JSON(w, r, render.M{
		"proxies":    proxies,
		"selectedId": selectedId,
	})
}

func addProxy(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	id, err := db2.AddProxy(req)
	if err != nil {
		log.Warnln("failed to add proxy: %v", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	render.JSON(w, r, render.M{"id": id})
}

func deleteProxy(w http.ResponseWriter, r *http.Request) {
	proxyId := chi.URLParam(r, "proxyId")
	err := db2.DeleteProxy(proxyId)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	render.NoContent(w, r)
}

func updateProxy(w http.ResponseWriter, r *http.Request) {
	proxyId := chi.URLParam(r, "proxyId")
	var req interface{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	if err := db2.UpdateProxy(proxyId, req.(map[string]interface{})); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	render.NoContent(w, r)
}
