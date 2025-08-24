package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/igoogolx/itun2socks/pkg/is_elevated"
)

func isAdminRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getIsAdmin)
	return r
}

var IsElevated = is_elevated.Get()

func getIsAdmin(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, render.M{
		"isAdmin": IsElevated,
	})
}
