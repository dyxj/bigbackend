package app

import (
	"net/http"

	"github.com/dyxj/bigbackend/pkg/idempotency"
)

func (s *Server) BuildRouter() http.Handler {
	router := http.NewServeMux()

	userProfileCreatorHandler, userProfileGetterHandler := s.buildUserProfileHandlers()

	router.HandleFunc("GET /healthz", s.healthCheck)
	router.Handle("GET /user/{id}/profile", s.TimeoutHandler(userProfileGetterHandler))
	router.Handle("POST /user/{id}/profile", s.TimeoutHandler(userProfileCreatorHandler))

	idemStore := idempotency.NewMemStore(idempotency.DefaultLockConfig)

	return idempotency.Middleware(
		s.logger,
		idempotency.DefaultKeyExtractor,
		idemStore,
		router,
	)
}
