package routes

import (
	"embed"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"io/fs"
	"net/http"
	"os"
)

type RouterHandler func(r chi.Router)

func Start(addr string) error {
	r := chi.NewRouter()
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
		r.Mount("/ping", pingRouter())
		r.Mount("/manager", managerRouter())
		r.Mount("/is-admin", isAdminRouter())
	})
	FileServer(r)
	err := http.ListenAndServe(addr, r)
	return err
}

//go:embed dist
var dashboard embed.FS

func FileServer(router *chi.Mux) {
	fSys, _ := fs.Sub(dashboard, "dist")

	staticFs := http.FileServer(http.FS(fSys))

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fSys.Open(r.RequestURI); os.IsNotExist(err) {
			http.StripPrefix(r.RequestURI, staticFs).ServeHTTP(w, r)
		} else {
			staticFs.ServeHTTP(w, r)
		}
	})
}
