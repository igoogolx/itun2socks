package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/igoogolx/itun2socks/global"
	"net/http"
)

func testRuleRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", testRule)
	return r
}

func testRule(w http.ResponseWriter, r *http.Request) {
	destination := r.URL.Query().Get("destination")
	if destination == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	if global.GetMatcher() == nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError("invalid matcher"))
		return
	}
	rule := global.GetMatcher().GetRule(destination)
	render.JSON(w, r, render.M{
		"rule": rule,
	})
}
