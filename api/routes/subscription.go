package routes

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/igoogolx/itun2socks/internal/cfg/outbound"
	"github.com/igoogolx/itun2socks/internal/configuration"
	"github.com/igoogolx/itun2socks/internal/conn"
	"github.com/igoogolx/itun2socks/internal/manager"
	"github.com/igoogolx/itun2socks/pkg/log"
)

func subscriptionRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/all", getSubscriptions)
	r.Put("/url", addSubscription)
	r.Delete("/", deleteSubscription)
	r.Post("/", updateSubscription)
	r.Post("/proxies", updateSubscriptionProxies)
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
