package idempotency

import (
	"context"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

const DefaultHeaderKey = "Idempotency-Key"
const TenantIdHeaderKey = "X-Tenant-Id"
const keySeparator = "::"

type Middleware struct {
	logger *zap.Logger
	store  Store
	config *Config
}

func NewMiddleware(
	logger *zap.Logger,
	store Store,
	options ...Option,
) *Middleware {
	config := DefaultConfig()

	for _, opt := range options {
		opt(config)
	}

	return &Middleware{
		logger: logger,
		store:  store,
		config: config,
	}
}

func (m *Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := m.config.extractor(r)
		if key == "" || r.Method == http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}

		err := m.store.Lock(r.Context(), key, m.config.lockOptions...)
		if err != nil {
			if errors.Is(err, ErrInProgress) {
				m.logger.Warn("idempotent request in progress", zap.String("key", key))
				m.config.errRespWriter(err, w)
				return
			}
			m.logger.Error("idempotent lock error", zap.Error(err), zap.String("key", key))
			m.config.errRespWriter(err, w)
			return
		}

		defer func(store Store, ctx context.Context, key string) {
			err := store.Unlock(ctx, key)
			if err != nil {
				m.logger.Error("failed to unlock idempotent key",
					zap.Error(err), zap.String("key", key))
			}
		}(m.store, r.Context(), key)

		cached, err := m.store.Get(r.Context(), key)
		if err != nil {
			m.logger.Error("failed to retrieve cache response", zap.Error(err), zap.String("key", key))
			m.config.errRespWriter(err, w)
			return
		}
		if cached != nil {
			serveCachedResponse(cached, w)
			return
		}

		// run next with response copier
		recorderWriter := newResponseRecorderWriter(w)

		// add idempotency key to request context and overwrite request
		ctx := WithValue(r.Context(), key)
		r = r.WithContext(ctx)

		next.ServeHTTP(recorderWriter, r)

		// cache response
		err = m.store.Set(ctx, key,
			&Response{
				Status: recorderWriter.status,
				Header: recorderWriter.cloneHeaders(),
				Body:   recorderWriter.body.Bytes(),
			},
			m.config.cacheExpiry)
		if err != nil {
			m.logger.Error("failed to store response", zap.Error(err), zap.String("key", key))
		}
	})
}

type ErrorResponseWriter func(err error, w http.ResponseWriter)

func DefaultErrorResponseWriter(err error, w http.ResponseWriter) {
	if errors.Is(err, ErrInProgress) {
		http.Error(w, "request with the same idempotency key is already in progress", http.StatusConflict)
		return
	}
	http.Error(w, "internal server error", http.StatusInternalServerError)
}

type KeyExtractor func(r *http.Request) string

func DefaultKeyExtractor(r *http.Request) string {
	return r.Header.Get(DefaultHeaderKey)
}

func TenantAndIdempotency(r *http.Request) string {
	key := r.Header.Get(DefaultHeaderKey)
	if key == "" {
		return ""
	}
	tenantId := r.Header.Get(TenantIdHeaderKey)
	return tenantId + keySeparator + key
}

func serveCachedResponse(cached *Response, w http.ResponseWriter) {
	clearAndCopyHeaders(w.Header(), cached.Header)
	cStatus := cached.Status
	if cStatus > 0 {
		// Only if status is set write it, else default behaviour would result in 200
		w.WriteHeader(cStatus)
	}
	// Errors are likely due to client disconnect.
	_, _ = w.Write(cached.Body)
}

func clearAndCopyHeaders(dst, src http.Header) {
	// ensure no prior headers leaks into response
	for k := range dst {
		dst.Del(k)
	}
	for k, vv := range src {
		copyV := make([]string, len(vv))
		copy(copyV, vv)
		dst[k] = copyV
	}
}
