package debug

import (
	"expvar"
	"net/http"
	"net/http/pprof"
)

func Mux() *http.ServeMux {
	m := http.NewServeMux()
	//register the debug handlers on this mux
	m.HandleFunc("/debug/pprof/", pprof.Index)
	m.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	m.HandleFunc("/debug/pprof/profile", pprof.Profile)
	m.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	m.HandleFunc("/debug/pprof/trace", pprof.Trace)

	//metrics
	m.Handle("/debug/vars/", expvar.Handler())
	return m
}
