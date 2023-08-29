package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	configuration2 "github.com/igoogolx/itun2socks/internal/configuration"
	"net/http"
)

func ruleRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getRules)
	r.Put("/", addRule)
	r.Delete("/{ruleId}", deleteRule)
	r.Post("/{ruleId}", updateRule)
	return r
}

func getRules(w http.ResponseWriter, r *http.Request) {
	rules, err := configuration2.GetRules()
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	selectId, err := configuration2.GetSelectedId("rule")
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

func addRule(w http.ResponseWriter, r *http.Request) {
	var req configuration2.RuleCfg
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	if _, err := configuration2.AddRule(req); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	render.NoContent(w, r)
}

func updateRule(w http.ResponseWriter, r *http.Request) {
	ruleId := chi.URLParam(r, "ruleId")
	var req configuration2.RuleCfg
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	if err := configuration2.UpdateRule(ruleId, req); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	render.NoContent(w, r)
}
func deleteRule(w http.ResponseWriter, r *http.Request) {
	ruleId := chi.URLParam(r, "ruleId")
	err := configuration2.DeleteRule(ruleId)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	render.NoContent(w, r)
}
