package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/gorilla/websocket"
	"github.com/igoogolx/itun2socks/internal/dns"
	metaDns "github.com/metacubex/mihomo/dns"
)

func dnsRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/statistic", getDnsStatistic)
	r.Post("/validate", validate)
	return r
}

type Dns struct {
	Success int32 `json:"success"`
	Fail    int32 `json:"fail"`
}

func getDnsStatistic(w http.ResponseWriter, r *http.Request) {
	var wsConn *websocket.Conn
	if websocket.IsWebSocketUpgrade(r) {
		var err error
		wsConn, err = upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
	}

	if wsConn == nil {
		w.Header().Set("Content-Type", "application/json")
		render.Status(r, http.StatusOK)
	}

	tick := time.NewTicker(time.Second)
	defer tick.Stop()
	buf := &bytes.Buffer{}
	var err error
	for range tick.C {
		buf.Reset()
		if err := json.NewEncoder(buf).Encode(dns.GetStatistic()); err != nil {
			break
		}

		if wsConn == nil {
			_, err = w.Write(buf.Bytes())
			w.(http.Flusher).Flush()
		} else {
			err = wsConn.WriteMessage(websocket.TextMessage, buf.Bytes())
		}

		if err != nil {
			break
		}
	}
}

type validateReq struct {
	servers []string
}

func validate(w http.ResponseWriter, r *http.Request) {

	var req validateReq
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}

	_, err := metaDns.ParseNameServer(req.servers)

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	render.NoContent(w, r)

}
