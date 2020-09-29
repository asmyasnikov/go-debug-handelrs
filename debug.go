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
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/allocs", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.HandleFunc("/debug/pprof/trace", pprof.Trace)

	router.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	router.Handle("/debug/pprof/block", pprof.Handler("block"))

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
