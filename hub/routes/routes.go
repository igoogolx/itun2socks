package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"net/http"
	"os"
)

type RouterHandler func(r chi.Router)

func Start(addr string, webDir string) error {
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
	FileServer(r, webDir)
	err := http.ListenAndServe(addr, r)
	return err
}

func FileServer(router *chi.Mux, webDir string) {
	fs := http.FileServer(http.Dir(webDir))

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(webDir + r.RequestURI); os.IsNotExist(err) {
			http.StripPrefix(r.RequestURI, fs).ServeHTTP(w, r)
		} else {
			fs.ServeHTTP(w, r)
		}
	})
}
