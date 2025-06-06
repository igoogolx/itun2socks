package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	configuration2 "github.com/igoogolx/itun2socks/internal/configuration"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

func settingRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getSetting)
	r.Get("/interfaces", getInterfaces)
	r.Get("/config-file-dir-path", getConfigDirPath)
	r.Get("/executable-path", getExecutablePath)
	r.Put("/", setSetting)
	r.Put("/reset-config", resetConfig)
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
	setting, err := configuration2.GetSetting()
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
	var req configuration2.SettingCfg
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	err := configuration2.SetSetting(req)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	render.NoContent(w, r)
}

func getConfigDirPath(w http.ResponseWriter, r *http.Request) {
	configFilePath, err := configuration2.GetConfigFilePath()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	dirPath := filepath.Dir(configFilePath)
	render.JSON(w, r, render.M{
		"path": dirPath,
	})
}

func getExecutablePath(w http.ResponseWriter, r *http.Request) {
	executablePath, err := os.Executable()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	render.JSON(w, r, render.M{
		"path": executablePath,
	})
}

func resetConfig(w http.ResponseWriter, r *http.Request) {
	err := configuration2.Reset()
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	render.NoContent(w, r)
}
