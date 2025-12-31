package app

import (
	"net/http"
)

func (s *Server) BuildRouter() http.Handler {
	router := http.NewServeMux()

	userProfileCreatorHandler, userProfileGetterHandler := s.buildUserProfileHandlers()

	router.HandleFunc("GET /healthz", s.healthCheck)
	router.Handle("GET /user/{id}/profile", s.TimeoutHandler(userProfileGetterHandler))
	router.Handle("POST /user/{id}/profile", s.TimeoutHandler(userProfileCreatorHandler))

	return router
}
