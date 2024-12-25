package routes

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/gorilla/websocket"
	"github.com/igoogolx/itun2socks/internal/manager"
	"net/http"
	"time"
)

func heartbeatRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/ping", ping)
	r.Get("/runtime-status", runtimeStatus)
	return r
}

func ping(w http.ResponseWriter, r *http.Request) {
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
	var err error
	for range tick.C {
		msg := []byte("pong")
		if wsConn == nil {
			_, err = w.Write(msg)
			w.(http.Flusher).Flush()
		} else {
			err = wsConn.WriteMessage(websocket.TextMessage, msg)
		}
		if err != nil {
			break
		}
	}
}

func runtimeStatus(w http.ResponseWriter, r *http.Request) {
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
	for range tick.C {
		msg, err := getRuntimeStatus()
		if err != nil {
			break
		}
		if wsConn == nil {
			msgData, err := json.Marshal(msg)
			if err != nil {
				break
			}
			_, err = w.Write(msgData)
			w.(http.Flusher).Flush()
		} else {
			err = wsConn.WriteJSON(msg)
		}
		if err != nil {
			break
		}
	}
}

type RuntimeStatus struct {
	Name      string `json:"name"`
	Addr      string `json:"addr"`
	IsStarted bool   `json:"isStarted"`
}

func getRuntimeStatus() (*RuntimeStatus, error) {
	isStarted := manager.GetIsStarted()
	name, addr := getCurProxy()
	return &RuntimeStatus{name, addr, isStarted}, nil
}
