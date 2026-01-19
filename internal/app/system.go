package app

import (
	"net/http"

	"github.com/dyxj/bigbackend/pkg/httpx"
)

func (s *Server) TimeoutHandler(h http.Handler) http.Handler {
	if s.httpConfig.HandlerTimeout() <= 0 {
		return h
	}
	return httpx.TimeoutHandler(h, s.httpConfig.HandlerTimeout())
}
