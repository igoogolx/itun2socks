package routes

import (
	"crypto/subtle"
	"embed"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/gorilla/websocket"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"unsafe"
)

var (
	serverSecret = ""
)

type RouterHandler func(r chi.Router)

func safeEuqal(a, b string) bool {
	aBuf := unsafe.Slice(unsafe.StringData(a), len(a))
	bBuf := unsafe.Slice(unsafe.StringData(b), len(b))
	return subtle.ConstantTimeCompare(aBuf, bBuf) == 1
}

func authentication(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if serverSecret == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Browser websocket not support custom header
		if websocket.IsWebSocketUpgrade(r) && r.URL.Query().Get("token") != "" {
			token := r.URL.Query().Get("token")
			if !safeEuqal(token, serverSecret) {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, ErrUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return
		}

		header := r.Header.Get("Authorization")
		bearer, token, found := strings.Cut(header, " ")

		hasInvalidHeader := bearer != "Bearer"
		hasInvalidSecret := !found || !safeEuqal(token, serverSecret)
		if hasInvalidHeader || hasInvalidSecret {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, ErrUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func Start(addr string, secret string) error {
	serverSecret = secret
	r := chi.NewRouter()
	r.Use(authentication)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*", "ws://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(middleware.Heartbeat("/ping"))
	r.Mount("/debug", middleware.Profiler())
	r.Group(func(r chi.Router) {
		r.Mount("/traffic", trafficRouter())
		r.Mount("/proxies", proxyRouter())
		r.Mount("/rules", ruleRouter())
		r.Mount("/selected", selectedRouter())
		r.Mount("/connection", connectionRouter())
		r.Mount("/log", logRouter())
		r.Mount("/setting", settingRouter())
		r.Mount("/version", versionRouter())
		r.Mount("/runtime-detail", runtimeDetailRouter())
		r.Mount("/manager", managerRouter())
		r.Mount("/is-admin", isAdminRouter())
		r.Mount("/heartbeat", heartbeatRouter())
	})
	go FileServer(r)
	err := http.ListenAndServe(addr, r)
	return err
}

//go:embed dist-ui
var dashboard embed.FS

func FileServer(router *chi.Mux) {
	fSys, _ := fs.Sub(dashboard, "dist-ui")

	staticFs := http.FileServer(http.FS(fSys))

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fSys.Open(r.RequestURI); os.IsNotExist(err) {
			http.StripPrefix(r.RequestURI, staticFs).ServeHTTP(w, r)
		} else {
			staticFs.ServeHTTP(w, r)
		}
	})
}
