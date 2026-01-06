package app

import (
	"errors"
	"net/http"
	"time"

	"github.com/dyxj/bigbackend/pkg/httpx"
	"github.com/dyxj/bigbackend/pkg/idempotency"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) BuildRouter() http.Handler {
	router := chi.NewRouter()

	idemStore := idempotency.NewMemStore(idempotency.DefaultLockConfig)
	idemMiddleware := idempotency.NewMiddleware(s.logger, idemStore,
		idempotency.WithCacheExpiry(24*time.Hour),
		idempotency.WithLockOptions(
			idempotency.WithLockRetry(3, 100*time.Millisecond),
			idempotency.WithLockExpiry(5*time.Second),
		),
		idempotency.WithErrorResponseWriter(s.idempotencyErrResponseWriter),
	)

	userProfileCreatorHandler, userProfileGetterHandler := s.buildUserProfileHandlers()

	router.Use(middleware.Logger)
	router.Get("/healthz", s.healthCheck)

	domainRouter := chi.NewRouter()

	domainRouter.Use(s.TimeoutHandler)
	domainRouter.Use(idemMiddleware.Handler)
	domainRouter.Use(middleware.Recoverer)

	domainRouter.Get("/user/{id}/profile", userProfileGetterHandler.ServeHTTP)
	domainRouter.Post("/user/{id}/profile", userProfileCreatorHandler.ServeHTTP)

	router.Mount("/", domainRouter)

	return router
}

func (s *Server) idempotencyErrResponseWriter(err error, w http.ResponseWriter) {
	if errors.Is(err, idempotency.ErrInProgress) {
		httpx.JsonResponse(
			http.StatusConflict,
			httpx.ErrorResponse{
				Code:    httpx.CodeIdempotencyError,
				Message: "processing of idempotency key is in progress",
			},
			w,
		)
	}

	httpx.InternalServerErrorResponse("", w)
}
