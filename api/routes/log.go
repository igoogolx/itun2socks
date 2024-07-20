package routes

import (
	"bytes"
	"encoding/json"
	"github.com/Dreamacro/clash/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"
	"github.com/igoogolx/itun2socks/internal/constants"
	"net/http"
	"strconv"
	"sync"
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
	var mux sync.Mutex
	logs := make([]Log, 0)
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
	go func() {
		for elm := range sub {
			func() {
				mux.Lock()
				defer mux.Unlock()
				buf.Reset()
				logEvent, ok := elm.(log.Event)
				if !ok {
					return
				}
				if logEvent.LogLevel < level {
					return
				}
				uid, _ := uuid.NewV4()
				logs = append(logs, Log{
					UUID:    uid,
					Time:    time.Now().UnixNano() / int64(time.Millisecond),
					Type:    logEvent.Type(),
					Payload: logEvent.Payload,
				})
			}()
		}
	}()

	sendLogs := func() error {
		mux.Lock()
		defer mux.Unlock()
		if err := json.NewEncoder(buf).Encode(logs); err != nil {
			return err
		}

		logs = make([]Log, 0)

		if wsConn == nil {
			_, err = w.Write(buf.Bytes())
			w.(http.Flusher).Flush()
		} else {
			err = wsConn.WriteMessage(websocket.TextMessage, buf.Bytes())
		}

		if err != nil {
			return err
		}
		return nil
	}

	intervalStr := r.URL.Query().Get("interval")
	interval := 1000
	if intervalStr != "" {
		t, err := strconv.Atoi(intervalStr)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, ErrBadRequest)
			return
		}

		interval = t
	}

	tick := time.NewTicker(time.Millisecond * time.Duration(interval))
	defer tick.Stop()
	for range tick.C {
		if err := sendLogs(); err != nil {
			break
		}
	}
}

func getLogsDir(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, render.M{
		"path": constants.Path.LogFilePath(),
	})
}
