package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/igoogolx/itun2socks/manager"
	"net/http"
)

func runtimeDetailRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getDetail)
	return r
}

func getDetail(w http.ResponseWriter, r *http.Request) {
	if !manager.GetIsStarted() {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError("the client is not started"))
		return
	}
	render.JSON(w, r, manager.RuntimeDetail())
}
