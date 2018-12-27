package serverpool

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/arjanvaneersel/kit/logger"
	"github.com/gorilla/mux"
)

type DiagnosticsServer struct {
	http.Server
	Router *mux.Router
	ready  *atomic.Value
	logger logger.Logger
}

func (s *DiagnosticsServer) Ready(v bool) {
	s.ready.Store(v)
}

func (s *DiagnosticsServer) Start() error {
	return s.ListenAndServe()
}

func (s *DiagnosticsServer) Stop(ctx context.Context) error {
	return s.Shutdown(ctx)
}

func (s *DiagnosticsServer) Address() string {
	return s.Addr
}

func (s *DiagnosticsServer) healthzHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *DiagnosticsServer) readyzHandler(w http.ResponseWriter, _ *http.Request) {
	if s.ready == nil || !s.ready.Load().(bool) {
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func NewDiagnosticsServer(port int, l logger.Logger) *DiagnosticsServer {
	s := DiagnosticsServer{
		ready: &atomic.Value{},
		Server: http.Server{
			Addr:           fmt.Sprintf(":%d", port),
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		logger: l,
	}

	r := mux.NewRouter()
	r.HandleFunc("/healthz", s.healthzHandler)
	r.HandleFunc("/readyz", s.readyzHandler)
	s.Handler = r

	s.ready.Store(false)
	return &s
}
