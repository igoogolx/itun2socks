package routes

import (
	"github.com/Dreamacro/clash/adapter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	configuration2 "github.com/igoogolx/itun2socks/internal/configuration"
	"github.com/igoogolx/itun2socks/internal/conn"
	"github.com/igoogolx/itun2socks/internal/executor"
	"github.com/igoogolx/itun2socks/internal/manager"
	"github.com/igoogolx/itun2socks/pkg/log"
	"net/http"
)

func selectedRouter() http.Handler {
	r := chi.NewRouter()
	r.Post("/rule", setRuleSelectedId)
	r.Post("/proxy", setProxySelectedId)
	return r
}

func setRuleSelectedId(w http.ResponseWriter, r *http.Request) {

	var req map[string]string
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	ruleSelectedId := req["id"]
	if err := configuration2.SetSelectedId("rule", ruleSelectedId); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	if manager.GetIsStarted() {
		ruleName, err := executor.UpdateRule()
		log.Infoln(log.FormatLog(log.ExecutorPrefix, "Update rule: %v"), ruleName)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, NewError(err.Error()))
			return
		}
	}

	render.NoContent(w, r)
}

func setProxySelectedId(w http.ResponseWriter, r *http.Request) {
	var req map[string]string
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return

	}
	proxySelectedId := req["id"]
	if err := configuration2.SetSelectedId("proxy", proxySelectedId); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	if manager.GetIsStarted() {
		rawProxy, err := configuration2.GetSelectedProxy()
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, ErrBadRequest)
			return
		}
		proxy, err := adapter.ParseProxy(rawProxy)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, ErrBadRequest)
			return
		}
		conn.UpdateProxy(proxy)
	}

	render.NoContent(w, r)
}
