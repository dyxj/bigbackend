package app

import (
	"errors"
	"net/http"
	"time"

	"github.com/dyxj/bigbackend/pkg/httpx"
	"github.com/dyxj/bigbackend/pkg/idempotency"
	"github.com/dyxj/bigbackend/pkg/monitoring"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) BuildRouter() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	if s.metrics != nil {
		router.Use(s.metrics.HTTPMetricsMiddleware)
	}

	router.Get("/healthz", monitoring.HealthCheckHandler(func() bool {
		return s.isShuttingDown.Load()
	}))
	router.Get("/readyz", monitoring.ReadinessCheckHandler(
		func() bool { return s.isShuttingDown.Load() },
		s.dbConn.Ping,
	))

	idemStore := idempotency.NewMemStore(idempotency.DefaultLockConfig)
	idemMiddleware := idempotency.NewMiddleware(s.logger, idemStore,
		idempotency.WithCacheExpiry(24*time.Hour),
		idempotency.WithLockOptions(
			idempotency.WithLockRetry(3, 100*time.Millisecond),
			idempotency.WithLockExpiry(5*time.Second),
		),
		idempotency.WithErrorResponseWriter(s.idempotencyErrResponseWriter),
	)

	apiRouter := chi.NewRouter()

	apiRouter.Use(s.TimeoutHandler)
	apiRouter.Use(idemMiddleware.Handler)
	apiRouter.Use(middleware.Recoverer)

	userProfileCreatorHandler, userProfileGetterHandler := s.buildUserProfileHandlers()
	apiRouter.Get("/user/{id}/profile", userProfileGetterHandler.ServeHTTP)
	apiRouter.Post("/user/{id}/profile", userProfileCreatorHandler.ServeHTTP)

	router.Mount("/api/v1", apiRouter)

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
