package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/igoogolx/itun2socks/internal/constants"
	"net/http"
)

func versionRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getVersion)
	return r
}

func getVersion(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, render.M{
		"version": constants.Version,
	})
}
