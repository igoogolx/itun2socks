package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/internal/manager"
	"net/http"
	"runtime"
)

type defaultDetail struct {
	HubAddress string `json:"hubAddress"`
}

func runtimeDetailRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getDetail)
	r.Get("/os", getOs)
	return r
}

func getOs(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, render.M{
		"os": runtime.GOOS,
	})
}

func getDetail(w http.ResponseWriter, r *http.Request) {
	hubAddress := constants.HubAddress()
	if !manager.GetIsStarted() {
		render.JSON(w, r, defaultDetail{
			HubAddress: hubAddress,
		})
		return
	}
	detail, err := manager.RuntimeDetail(hubAddress)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, err)
		return
	}
	render.JSON(w, r, detail)
}
