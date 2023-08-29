package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/igoogolx/itun2socks/pkg/is_elevated"
	"net/http"
)

func isAdminRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getIsAdmin)
	return r
}

func getIsAdmin(w http.ResponseWriter, r *http.Request) {
	isAdmin := is_elevated.Get()
	render.JSON(w, r, render.M{
		"isAdmin": isAdmin,
	})
}
