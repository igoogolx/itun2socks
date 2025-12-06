package routes

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/Dreamacro/clash/adapter"
	C "github.com/Dreamacro/clash/constant"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/igoogolx/itun2socks/internal/cfg/outbound"
	"github.com/igoogolx/itun2socks/internal/configuration"
	"github.com/igoogolx/itun2socks/internal/conn"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/internal/manager"
	"github.com/igoogolx/itun2socks/internal/tunnel"
	"github.com/igoogolx/itun2socks/pkg/log"
)

var (
	defaultDelayTimeout = 5 * time.Second
	defaultDelayTestUrl = "https://www.google.com"
)

func proxyRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getProxies)
	r.Get("/cur-proxy", handleGetProxy)
	r.Post("/url", getResFromUrl)
	r.Put("/", addProxy)
	r.Delete("/all", deleteAllProxies)
	r.Delete("/", deleteProxies)
	r.Post("/{proxyId}", updateProxy)
	r.Get("/delay/{proxyId}", getProxyDelay)
	r.Get("/udp-test/{proxyId}", testProxyUdp)
	r.Get("/subscriptions", getSubscriptions)
	r.Put("/subscription-url", addSubscription)
	r.Delete("/subscription", deleteSubscription)
	r.Post("/subscription", updateSubscription)
	r.Post("/subscription/proxies", updateSubscriptionProxies)
	return r
}

func getResFromUrl(w http.ResponseWriter, r *http.Request) {
	var reqParams map[string]string
	if err := render.DecodeJSON(r.Body, &reqParams); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	client := &http.Client{}

	req, err := http.NewRequest("GET", reqParams["url"], nil)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36 Edg/141.0.0.0")

	resp, err := client.Do(req)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Debugln("%s", log.FormatLog(log.HubPrefix, "get res from url: fail to close body"))
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}

	render.JSON(w, r, render.M{
		"data": string(body),
	})
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
	proxyOption, err := configuration.GetProxy(proxyId)
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
	metadata, err := tunnel.CreateMetadata("0.0.0.0:0", "8.8.8.8:53", C.UDP)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	pc, err := p.ListenPacketContext(context.Background(), metadata)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	res, err := UdpTest(pc, "8.8.8.8:53")
	if err != nil {
		log.Warnln(log.FormatLog(log.HubPrefix, "fail to test udp, err: %v"), err)
		res = false
	}

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
	proxyOption, err := configuration.GetProxy(proxyId)
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
	delay, _, err := p.URLTest(ctx, url)
	if err != nil {
		render.JSON(w, r, render.M{
			"delay": -1,
		})
		return
	}
	proxyOption["delay"] = delay
	err = configuration.UpdateProxy(proxyId, proxyOption)
	if err != nil {
		render.JSON(w, r, render.M{
			"delay": -1,
		})
		return
	}
	render.JSON(w, r, render.M{
		"delay": delay,
	})
}

func getProxies(w http.ResponseWriter, r *http.Request) {
	proxiesMap, err := configuration.GetProxies()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	selectedId, err := configuration.GetSelectedId("proxy")
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

func getCurProxy() (string, string) {
	name := ""
	addr := ""

	if manager.GetIsStarted() {
		curAutoProxy, err := conn.GetProxy(constants.PolicyProxy)
		if err == nil {
			if curAutoProxy.Type() == C.URLTest || curAutoProxy.Type() == C.Fallback {
				curAutoProxy = curAutoProxy.Unwrap(&C.Metadata{})
			}
		}
		if curAutoProxy != nil {
			name = curAutoProxy.Name()
			addr = curAutoProxy.Addr()
		}
	} else {
		curSelectedProxy, err := configuration.GetSelectedProxy()
		if err == nil {
			if proxyName, ok := curSelectedProxy["name"].(string); ok {
				name = proxyName
			}
			if proxyAddr, ok := curSelectedProxy["server"].(string); ok {
				addr = proxyAddr
			}
		}
	}

	return name, addr

}

func handleGetProxy(w http.ResponseWriter, r *http.Request) {
	name, addr := getCurProxy()
	render.JSON(w, r, render.M{
		"name": name,
		"addr": addr,
	})
}

func addProxy(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	id, err := configuration.AddProxy(req)
	if err != nil {
		log.Warnln(log.FormatLog(log.HubPrefix, "fail to add proxy: %v"), err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	render.JSON(w, r, render.M{"id": id})
}

func deleteProxies(w http.ResponseWriter, r *http.Request) {
	var req map[string][]string
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	err := configuration.DeleteProxies(req["ids"])
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	render.NoContent(w, r)
}

func deleteAllProxies(w http.ResponseWriter, r *http.Request) {
	err := configuration.DeleteAllProxies()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	render.NoContent(w, r)
}

func addSubscription(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Proxies            []map[string]interface{} `json:"proxies"`
		SubscriptionUrl    string                   `json:"subscriptionUrl"`
		SubscriptionName   string                   `json:"subscriptionName"`
		SubscriptionRemark string                   `json:"subscriptionRemark"`
	}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	proxies := req.Proxies
	if proxies == nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError("invalid proxies"))
		return
	}
	newProxies, newSubscriptions, err := configuration.AddSubscription(proxies, req.SubscriptionUrl, req.SubscriptionName, req.SubscriptionRemark)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	rawConfig, err := configuration.Read()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	if manager.GetIsStarted() && rawConfig.Setting.AutoMode.Enabled {
		outboundOption := outbound.Option{
			AutoMode:      rawConfig.Setting.AutoMode,
			Proxies:       rawConfig.Proxy,
			SelectedProxy: rawConfig.Selected.Proxy,
		}
		proxy, err := outbound.New(outboundOption)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, NewError(err.Error()))
			return
		}
		conn.UpdateProxy(proxy)
	}
	render.JSON(w, r, render.M{"proxies": newProxies, "subscriptions": newSubscriptions})
}

func updateProxy(w http.ResponseWriter, r *http.Request) {
	proxyId := chi.URLParam(r, "proxyId")
	var req interface{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	if err := configuration.UpdateProxy(proxyId, req.(map[string]interface{})); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	render.NoContent(w, r)
}

func getSubscriptions(w http.ResponseWriter, r *http.Request) {
	subscriptions, err := configuration.GetSubscriptions()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	render.JSON(w, r, render.M{
		"subscriptions": subscriptions,
	})
}

func deleteSubscription(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Id string `json:"id"`
	}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	err := configuration.DeleteSubscription(req.Id)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	render.NoContent(w, r)
}

func updateSubscription(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Subscription configuration.SubscriptionCfg `json:"subscription"`
	}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	err := configuration.UpdateSubscription(req.Subscription)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	render.NoContent(w, r)
}

func updateSubscriptionProxies(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SubscriptionId string                   `json:"subscriptionId"`
		Proxies        []map[string]interface{} `json:"proxies"`
	}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	newProxies, err := configuration.UpdateSubscriptionProxies(req.SubscriptionId, req.Proxies)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}

	render.JSON(w, r, render.M{"proxies": newProxies})
}
