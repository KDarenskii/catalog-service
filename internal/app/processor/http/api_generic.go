package rprocessor

import (
	"net/http"
	"net/http/pprof"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	rhandler "github.com/KDarenskii/catalog-service/internal/app/handler/http"
)

func vGenericRegHealthCheck(r *mux.Router, h rhandler.Health) {
	reg(r, http.MethodGet, "/health", http.HandlerFunc(h.LastCheck))
	reg(r, http.MethodGet, "/metrics", promhttp.Handler())
}

func handlerNotFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func vGenericRegPprof(r *mux.Router) {
	reg(r, http.MethodGet, "/debug/pprof", http.HandlerFunc(pprofIndex))
	reg(r, http.MethodGet, "/debug/pprof/", http.HandlerFunc(pprof.Index))

	reg(r, http.MethodGet, "/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	reg(r, http.MethodGet, "/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	reg(r, http.MethodGet, "/debug/pprof/trace", http.HandlerFunc(pprof.Trace))

	r.Methods(http.MethodGet, http.MethodPost).Path("/debug/pprof/symbol").Handler(http.HandlerFunc(pprof.Symbol))

	reg(r, http.MethodGet, "/debug/pprof/allocs", pprof.Handler("allocs"))
	reg(r, http.MethodGet, "/debug/pprof/block", pprof.Handler("block"))
	reg(r, http.MethodGet, "/debug/pprof/goroutine", pprof.Handler("goroutine"))
	reg(r, http.MethodGet, "/debug/pprof/heap", pprof.Handler("heap"))
	reg(r, http.MethodGet, "/debug/pprof/mutex", pprof.Handler("mutex"))
	reg(r, http.MethodGet, "/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
}

func pprofIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/debug/pprof/", http.StatusMovedPermanently)
}
