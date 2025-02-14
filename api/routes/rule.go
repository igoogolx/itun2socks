package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/igoogolx/itun2socks/internal/cfg/distribution/ruleEngine"
	"github.com/igoogolx/itun2socks/internal/configuration"
	"github.com/igoogolx/itun2socks/internal/executor"
	"github.com/igoogolx/itun2socks/internal/manager"
	"github.com/igoogolx/itun2socks/internal/tunnel/statistic"
	"net/http"
)

func ruleRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getRules)
	r.Get("/{id}", getRuleDetail)
	r.Put("/customized", addCustomizedRules)
	r.Post("/customized", editCustomizedRule)
	r.Delete("/customized", deleteCustomizedRules)
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

func addCustomizedRules(w http.ResponseWriter, r *http.Request) {
	var req map[string][]string
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	if len(req["rules"]) == 0 {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, NewError("invalid rules"))
		return
	}
	err := configuration.AddCustomizedRule(req["rules"])
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	if manager.GetIsStarted() {
		_, err := executor.UpdateRule()
		statistic.DefaultManager.CloseAllConnections()
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, NewError(err.Error()))
			return
		}
	}
	render.NoContent(w, r)
}

func deleteCustomizedRules(w http.ResponseWriter, r *http.Request) {
	var req map[string][]string
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	if len(req["rules"]) == 0 {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, NewError("invalid rules"))
		return
	}
	err := configuration.DeleteCustomizedRule(req["rules"])
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}

	if manager.GetIsStarted() {
		_, err := executor.UpdateRule()
		statistic.DefaultManager.CloseAllConnections()
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, NewError(err.Error()))
			return
		}
	}

	render.NoContent(w, r)
}

func getRuleDetail(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	var rules []ruleEngine.Rule
	var err error
	if id == "customized" {
		rules, err = configuration.GetCustomizedRules()
	} else {
		rules, err = configuration.GetBuiltInRules(id)
	}
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	render.JSON(w, r, render.M{
		"items": rules,
	})
}

func editCustomizedRule(w http.ResponseWriter, r *http.Request) {
	var req map[string]string
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	if len(req["oldRule"]) == 0 || len(req["newRule"]) == 0 {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, NewError("invalid rules"))
		return
	}
	err := configuration.EditCustomizedRule(req["oldRule"], req["newRule"])
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	if manager.GetIsStarted() {
		_, err := executor.UpdateRule()
		statistic.DefaultManager.CloseAllConnections()
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, NewError(err.Error()))
			return
		}
	}
	render.NoContent(w, r)
}
