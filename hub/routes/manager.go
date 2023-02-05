package routes

import (
	"github.com/Dreamacro/clash/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/igoogolx/itun2socks/manager"
	"net/http"
)

func managerRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getStatus)
	r.Post("/start", start)
	r.Post("/stop", stop)
	return r
}

func start(w http.ResponseWriter, r *http.Request) {
	err := manager.Start()
	if err != nil {
		log.Errorln("fail to start the client,error:%v", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	render.NoContent(w, r)
}

func stop(w http.ResponseWriter, r *http.Request) {
	err := manager.Close()
	if err != nil {
		log.Errorln("fail to stop the client,error:%v", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	render.NoContent(w, r)
}

type Status struct {
	IsStarted bool `json:"isStarted"`
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	isStarted := manager.GetIsStarted()
	status := &Status{isStarted}
	render.JSON(w, r, status)
}
