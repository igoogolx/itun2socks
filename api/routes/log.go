package routes

import (
	"bytes"
	"encoding/json"
	cLog "github.com/Dreamacro/clash/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/gorilla/websocket"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/pkg/log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Log struct {
	Type string `json:"level"`
	Time int64  `json:"time"`
	Msg  string `json:"msg"`
}

func (l Log) String() (string, error) {
	jsonBytes, err := json.Marshal(l)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func logRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getLogs)
	r.Get("/dir", getLogsDir)
	return r
}

func getLogs(w http.ResponseWriter, r *http.Request) {
	var mux sync.Mutex
	logs, err := log.ReadFile(1000)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, err)
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

	sub := cLog.Subscribe()
	defer cLog.UnSubscribe(sub)
	go func() {
		for elm := range sub {
			func() {
				mux.Lock()
				defer mux.Unlock()
				logEvent, ok := elm.(cLog.Event)
				if !ok {
					return
				}
				logValue, err := Log{
					Time: time.Now().UnixNano() / int64(time.Millisecond),
					Type: logEvent.Type(),
					Msg:  logEvent.Payload,
				}.String()
				if err != nil {
					return
				}
				logs = append(logs, logValue)
			}()
		}
	}()

	buf := &bytes.Buffer{}
	sendLogs := func() error {
		mux.Lock()
		defer mux.Unlock()

		buf.Reset()
		if err := json.NewEncoder(buf).Encode(logs); err != nil {
			return err
		}

		logs = make([]string, 0)

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
