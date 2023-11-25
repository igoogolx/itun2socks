package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/igoogolx/itun2socks/internal/configuration"
	"net/http"
)

func ruleRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getRules)
	r.Get("/{id}", getRuleDetail)
	return r
}

func getRules(w http.ResponseWriter, r *http.Request) {
	rules, err := configuration.GetRuleIds()
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	selectId, err := configuration.GetSelectedId("rule")
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	render.JSON(w, r, render.M{
		"rules":      rules,
		"selectedId": selectId,
	})
}

func getRuleDetail(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	rules, err := configuration.GetBuiltInRules(id)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	render.JSON(w, r, render.M{
		"items": rules,
	})
}
