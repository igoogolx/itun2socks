package routes

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/gorilla/websocket"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/internal/tunnel/statistic"
	"net/http"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func trafficRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getTraffic)
	return r
}

type TrafficItem struct {
	Up   int64 `json:"upload"`
	Down int64 `json:"download"`
}

type Speed struct {
	Proxy  TrafficItem `json:"proxy"`
	Direct TrafficItem `json:"direct"`
}

type Traffic struct {
	Speed Speed           `json:"speed"`
	Total statistic.Total `json:"total"`
}

func getTraffic(w http.ResponseWriter, r *http.Request) {
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
	t := statistic.DefaultManager
	buf := &bytes.Buffer{}
	var err error
	for range tick.C {
		buf.Reset()
		proxyUp, proxyDown := t.Now(constants.PolicyProxy)
		directUp, directDown := t.Now(constants.PolicyDirect)
		if err := json.NewEncoder(buf).Encode(Traffic{
			Speed: Speed{
				TrafficItem{proxyUp, proxyDown},
				TrafficItem{directUp, directDown},
			},
			Total: *t.GetTotal(),
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
