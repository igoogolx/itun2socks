package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

func heartbeatRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/ping", ping)
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
