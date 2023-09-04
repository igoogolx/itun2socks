package routes

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/igoogolx/itun2socks/internal/tunnel/statistic"
	"github.com/igoogolx/itun2socks/pkg/log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/render"
	"github.com/gorilla/websocket"
)

func connectionRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getConnections)
	r.Delete("/", closeAllConnections)
	r.Delete("/{id}", closeConnection)
	return r
}

func getConnections(w http.ResponseWriter, r *http.Request) {
	if !websocket.IsWebSocketUpgrade(r) {
		snapshot := statistic.DefaultManager.Connections()
		render.JSON(w, r, snapshot)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
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

	buf := &bytes.Buffer{}
	sendSnapshot := func() error {
		buf.Reset()
		snapshot := statistic.DefaultManager.Connections()
		if err := json.NewEncoder(buf).Encode(snapshot); err != nil {
			return err
		}

		return conn.WriteMessage(websocket.TextMessage, buf.Bytes())
	}

	if err := sendSnapshot(); err != nil {
		return
	}

	tick := time.NewTicker(time.Millisecond * time.Duration(interval))
	defer tick.Stop()
	for range tick.C {
		if err := sendSnapshot(); err != nil {
			break
		}
	}
}

func closeConnection(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	connections := statistic.DefaultManager.Connections()
	for _, c := range connections {
		if id == c.ID() {
			err := c.Close()
			if err != nil {
				log.Debugln(log.FormatLog(log.HubPrefix, "fail to close connection, err: %v"), err)
			}
			break
		}
	}
	render.NoContent(w, r)
}

func closeAllConnections(w http.ResponseWriter, r *http.Request) {
	statistic.DefaultManager.CloseAllConnections()
	render.NoContent(w, r)
}
