package app

import (
	"errors"
	"net/http"
	"time"

	"github.com/dyxj/bigbackend/pkg/httpx"
	"github.com/dyxj/bigbackend/pkg/idempotency"
)

func (s *Server) BuildRouter() http.Handler {
	router := http.NewServeMux()

	userProfileCreatorHandler, userProfileGetterHandler := s.buildUserProfileHandlers()

	router.HandleFunc("GET /healthz", s.healthCheck)
	router.Handle("GET /user/{id}/profile", s.TimeoutHandler(userProfileGetterHandler))
	router.Handle("POST /user/{id}/profile", s.TimeoutHandler(userProfileCreatorHandler))

	idemStore := idempotency.NewMemStore(idempotency.DefaultLockConfig)
	idemMiddleware := idempotency.NewMiddleware(s.logger, idemStore,
		idempotency.WithCacheExpiry(24*time.Hour),
		idempotency.WithLockOptions(
			idempotency.WithLockRetry(3, 100*time.Millisecond),
			idempotency.WithLockExpiry(5*time.Second),
		),
		idempotency.WithErrorResponseWriter(s.idempotencyErrResponseWriter),
	)

	return idemMiddleware.Handler(router)
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
