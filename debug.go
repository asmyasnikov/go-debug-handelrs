package debug

import (
	"github.com/dustin/go-humanize"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/http/pprof"
	"runtime"
	"strconv"
	"time"
)

// Serve make handlers for debug profiling
func Serve(port int) *http.Server {
	if port <= 1024 {
		return nil
	}
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: Handlers(mux.NewRouter()),
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Error().Caller().Err(err).Msg("Error on start web-server")
		}
	}()
	return server
}

// Handlers register debug handlers on router
func Handlers(router *mux.Router) *mux.Router {
	r := router.PathPrefix("/debug/pprof/").Subrouter()
	r.HandleFunc("/", pprof.Index)
	r.HandleFunc("/allocs", pprof.Cmdline)
	r.HandleFunc("/cmdline", pprof.Cmdline)
	r.HandleFunc("/profile", pprof.Profile)
	r.HandleFunc("/symbol", pprof.Symbol)
	r.HandleFunc("/trace", pprof.Trace)

	r.Handle("/mutex", pprof.Handler("mutex"))
	r.Handle("/goroutine", pprof.Handler("goroutine"))
	r.Handle("/heap", pprof.Handler("heap"))
	r.Handle("/threadcreate", pprof.Handler("threadcreate"))
	r.Handle("/block", pprof.Handler("block"))

	go func() {
		var memory runtime.MemStats
		for {
			runtime.ReadMemStats(&memory)
			log.
				Info().
				Caller().
				Str("stack", humanize.Bytes(memory.StackInuse)).
				Str("heap", humanize.Bytes(memory.HeapAlloc)).
				Str("total", humanize.Bytes(memory.StackInuse+memory.HeapAlloc)).
				Str("sys", humanize.Bytes(memory.Sys)).
				Uint32("num gc", memory.NumGC).
				Msg("")
			time.Sleep(time.Second)
		}
	}()

	return router
}
