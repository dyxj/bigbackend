package app

import (
	"net/http"

	"github.com/dyxj/bigbackend/pkg/httpx"
	"go.uber.org/zap"
)

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	if s.isShuttingDown.Load() {
		http.Error(w, "shutting down", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		s.logger.Warn("Error writing health check response", zap.Error(err))
		return
	}
}

func (s *Server) TimeoutHandler(h http.Handler) http.Handler {
	if s.httpConfig.HandlerTimeout() <= 0 {
		return h
	}
	return httpx.TimeoutHandler(h, s.httpConfig.HandlerTimeout())
}
