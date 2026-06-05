package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/igoogolx/itun2socks/internal/configuration"
	"github.com/igoogolx/itun2socks/internal/conn"
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/internal/manager"
	"github.com/igoogolx/itun2socks/internal/tunnel"
	"github.com/igoogolx/itun2socks/pkg/clash/adapter"
	C "github.com/igoogolx/itun2socks/pkg/clash/constant"
	"github.com/igoogolx/itun2socks/pkg/log"
)

var (
	defaultDelayTimeout = 5 * time.Second
	defaultDelayTestUrl = "https://www.google.com"
)

func proxyRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getProxies)
	r.Get("/cur-proxy", handleGetProxy)
	r.Get("/{proxyId}", getProxy) // NEW: get single proxy with password
	r.Put("/", addProxy)
	r.Delete("/all", deleteAllProxies)
	r.Delete("/", deleteProxies)
	r.Post("/{proxyId}/lock-password", lockProxyPassword)
	r.Post("/{proxyId}/reset-password", resetProxyPassword)
	r.Get("/{proxyId}/password-status", getPasswordStatus)
	r.Post("/{proxyId}", updateProxy)
	r.Get("/delay/{proxyId}", getProxyDelay)
	r.Get("/udp-test/{proxyId}", testProxyUdp)
	return r
}

func testProxyUdp(w http.ResponseWriter, r *http.Request) {
	proxyId := chi.URLParam(r, "proxyId")
	if proxyId == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	url := chi.URLParam(r, "url")
	if url == "" {
		url = defaultDelayTestUrl
	}
	proxyOption, err := configuration.GetProxy(proxyId)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	p, err := adapter.ParseProxy(proxyOption)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	metadata, err := tunnel.CreateMetadata("0.0.0.0:0", "8.8.8.8:53", C.UDP)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	pc, err := p.ListenPacketContext(context.Background(), metadata)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	res, err := UdpTest(pc, "8.8.8.8:53")
	if err != nil {
		log.Warnln(log.FormatLog(log.HubPrefix, "fail to test udp, err: %v"), err)
		res = false
	}

	render.JSON(w, r, render.M{
		"result": res,
	})
}

func getProxyDelay(w http.ResponseWriter, r *http.Request) {
	proxyId := chi.URLParam(r, "proxyId")
	if proxyId == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	url := chi.URLParam(r, "url")
	if url == "" {
		url = defaultDelayTestUrl
	}
	proxyOption, err := configuration.GetProxy(proxyId)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	p, err := adapter.ParseProxy(proxyOption)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultDelayTimeout)
	defer cancel()
	delay, _, err := p.URLTest(ctx, url)
	if err != nil {
		render.JSON(w, r, render.M{
			"delay": -1,
		})
		return
	}
	proxyOption["delay"] = delay
	err = configuration.UpdateProxy(proxyId, proxyOption)
	if err != nil {
		render.JSON(w, r, render.M{
			"delay": -1,
		})
		return
	}
	render.JSON(w, r, render.M{
		"delay": delay,
	})
}

// getProxies returns all proxies with passwords STRIPPED for security.
// The web dashboard cannot see passwords through this endpoint.
func getProxies(w http.ResponseWriter, r *http.Request) {
	proxiesMap, err := configuration.GetProxies()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	selectedId, err := configuration.GetSelectedId("proxy")
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	proxies := make([]any, 0)
	for _, proxy := range proxiesMap {
		// Strip passwords from the list — use getProxy/:id to retrieve them
		proxies = append(proxies, configuration.StripProxyPasswords(proxy))
	}
	render.JSON(w, r, render.M{
		"proxies":    proxies,
		"selectedId": selectedId,
	})
}

// getProxy returns a single proxy INCLUDING its password.
// If the proxy is locked, passwords are stripped from the response.
func getProxy(w http.ResponseWriter, r *http.Request) {
	proxyId := chi.URLParam(r, "proxyId")
	if proxyId == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	proxy, err := configuration.GetProxy(proxyId)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	// If password is locked, strip it from response
	if locked, ok := proxy["passwordLocked"].(bool); ok && locked {
		result := configuration.StripProxyPasswords(proxy)
		result["passwordLocked"] = true
		render.JSON(w, r, result)
		return
	}
	render.JSON(w, r, proxy)
}

func getCurProxy() (string, string) {
	name := ""
	addr := ""

	if manager.GetIsStarted() {
		curAutoProxy, err := conn.GetProxy(constants.PolicyProxy)
		if err == nil {
			if curAutoProxy.Type() == C.URLTest || curAutoProxy.Type() == C.Fallback {
				curAutoProxy = curAutoProxy.Unwrap(&C.Metadata{})
			}
		}
		if curAutoProxy != nil {
			name = curAutoProxy.Name()
			addr = curAutoProxy.Addr()
		}
	} else {
		curSelectedProxy, err := configuration.GetSelectedProxy()
		if err == nil {
			if proxyName, ok := curSelectedProxy["name"].(string); ok {
				name = proxyName
			}
			if proxyAddr, ok := curSelectedProxy["server"].(string); ok {
				addr = proxyAddr
			}
		}
	}

	return name, addr
}

func handleGetProxy(w http.ResponseWriter, r *http.Request) {
	name, addr := getCurProxy()
	render.JSON(w, r, render.M{
		"name": name,
		"addr": addr,
	})
}

func addProxy(w http.ResponseWriter, r *http.Request) {
	var req map[string]any
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	id, err := configuration.AddProxy(req)
	if err != nil {
		log.Warnln(log.FormatLog(log.HubPrefix, "fail to add proxy: %v"), err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	render.JSON(w, r, render.M{"id": id})
}

func deleteProxies(w http.ResponseWriter, r *http.Request) {
	var req map[string][]string
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	err := configuration.DeleteProxies(req["ids"])
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	render.NoContent(w, r)
}

func deleteAllProxies(w http.ResponseWriter, r *http.Request) {
	err := configuration.DeleteAllProxies()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	render.NoContent(w, r)
}

func updateProxy(w http.ResponseWriter, r *http.Request) {
	proxyId := chi.URLParam(r, "proxyId")
	var req any
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}

	// Auto stop/start if the manager is running (allows save while active)
	wasRunning := manager.GetIsStarted()
	if wasRunning {
		_ = manager.Close()
	}

	if err := configuration.UpdateProxy(proxyId, req.(map[string]any)); err != nil {
		// Restart if we stopped
		if wasRunning {
			_ = manager.Start()
		}
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}

	// Restart if it was running
	if wasRunning {
		_ = manager.Start()
	}

	render.NoContent(w, r)
}

func lockProxyPassword(w http.ResponseWriter, r *http.Request) {
	proxyId := chi.URLParam(r, "proxyId")
	if proxyId == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	if err := configuration.LockProxyPassword(proxyId); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	render.NoContent(w, r)
}

func resetProxyPassword(w http.ResponseWriter, r *http.Request) {
	proxyId := chi.URLParam(r, "proxyId")
	if proxyId == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	if err := configuration.ResetProxyPassword(proxyId); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, NewError(err.Error()))
		return
	}
	render.NoContent(w, r)
}

func getPasswordStatus(w http.ResponseWriter, r *http.Request) {
	proxyId := chi.URLParam(r, "proxyId")
	if proxyId == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
	proxy, err := configuration.GetProxy(proxyId)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, NewError(err.Error()))
		return
	}

	mode, _ := proxy["passwordMode"].(string)
	locked, _ := proxy["passwordLocked"].(bool)
	expired := configuration.CheckPasswordExpiry(proxy)
	hasPassword := false
	for _, field := range []string{"password", "passwd", "auth_str"} {
		if val, ok := proxy[field].(string); ok && val != "" {
			hasPassword = true
			break
		}
	}

	render.JSON(w, r, render.M{
		"mode":        mode,
		"locked":      locked,
		"expired":     expired,
		"hasPassword": hasPassword,
		"setAt":       proxy["passwordSetAt"],
		"ttlMinutes":  proxy["passwordTTLMinutes"],
	})
}
