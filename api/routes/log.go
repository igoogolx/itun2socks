package routes

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"
	"github.com/igoogolx/clash/log"
	"github.com/igoogolx/itun2socks/internal/constants"
	"net/http"
	"time"
)

type Log struct {
	UUID    uuid.UUID `json:"id"`
	Type    string    `json:"type"`
	Time    int64     `json:"time"`
	Payload string    `json:"payload"`
}

func logRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getLogs)
	r.Get("/dir", getLogsDir)
	return r
}

func getLogs(w http.ResponseWriter, r *http.Request) {
	levelText := r.URL.Query().Get("level")
	if levelText == "" {
		levelText = "info"
	}

	level, ok := log.LogLevelMapping[levelText]
	if !ok {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}

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

	sub := log.Subscribe()
	defer log.UnSubscribe(sub)
	buf := &bytes.Buffer{}
	var err error
	for elm := range sub {
		buf.Reset()
		logEvent, ok := elm.(log.Event)
		if !ok {
			break
		}
		if logEvent.LogLevel < level {
			continue
		}

		uid, _ := uuid.NewV4()
		if err := json.NewEncoder(buf).Encode(Log{
			UUID:    uid,
			Time:    time.Now().UnixNano() / int64(time.Millisecond),
			Type:    logEvent.Type(),
			Payload: logEvent.Payload,
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

func getLogsDir(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, render.M{
		"path": constants.Path.LogFilePath(),
	})
}
