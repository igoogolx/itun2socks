package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/igoogolx/itun2socks/configuration"
	"github.com/igoogolx/itun2socks/configuration/configuration-types"
	"net"
	"net/http"
)

func settingRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getSetting)
	r.Get("/interfaces", getInterfaces)
	r.Put("/", setSetting)
	return r
}

func getInterfaces(w http.ResponseWriter, r *http.Request) {
	interfaces, err := net.Interfaces()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	render.JSON(w, r, render.M{
		"interfaces": interfaces,
	})
}

func getSetting(w http.ResponseWriter, r *http.Request) {
	setting, err := configuration.GetSetting()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	render.JSON(w, r, render.M{
		"setting": setting,
	})
}

func setSetting(w http.ResponseWriter, r *http.Request) {
	var req configuration_types.SettingCfg
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	err := configuration.SetSetting(req)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	render.NoContent(w, r)
}
