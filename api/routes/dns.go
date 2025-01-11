package routes

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/gorilla/websocket"
	"github.com/igoogolx/itun2socks/internal/dns"
	"net/http"
	"time"
)

func dnsRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/statistic", getDnsStatistic)
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
		if err := json.NewEncoder(buf).Encode(Dns{
			Success: dns.GetSuccessQueryCount(),
			Fail:    dns.GetFailQueryCount(),
		}); err != nil {
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
