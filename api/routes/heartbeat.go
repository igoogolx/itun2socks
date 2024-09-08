package routes

import (
	"encoding/json"
	C "github.com/Dreamacro/clash/constant"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/gorilla/websocket"
	"github.com/igoogolx/itun2socks/internal/conn"
	"github.com/igoogolx/itun2socks/internal/constants"
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
	name := ""
	addr := ""
	isStarted := manager.GetIsStarted()
	if isStarted {
		curAutoProxy, err := conn.GetProxy(constants.PolicyProxy)
		if err != nil {
			return nil, err
		}
		if curAutoProxy != nil {
			if curAutoProxy.Type() == C.URLTest || curAutoProxy.Type() == C.Fallback {
				curAutoProxy = curAutoProxy.Unwrap(&C.Metadata{})
			}
		}
		if curAutoProxy != nil {
			name = curAutoProxy.Name()
			addr = curAutoProxy.Addr()
		}
	}

	return &RuntimeStatus{name, addr, isStarted}, nil

}
