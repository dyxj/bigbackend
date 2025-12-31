package idempotency

import (
	"context"
	"errors"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const DefaultHeaderKey = "Idempotency-Key"
const TenantIdHeaderKey = "X-Tenant-Id"
const keySeparator = "::"

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

func Middleware(
	logger *zap.Logger,
	extractor KeyExtractor,
	store Store,
	next http.Handler,
) http.Handler {
	if extractor == nil {
		extractor = DefaultKeyExtractor
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := extractor(r)
		if key == "" || r.Method == http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}

		err := store.Lock(r.Context(), key, WithRetry(2, 100*time.Millisecond))
		if err != nil {
			if errors.Is(err, ErrInProgress) {
				http.Error(w, "idempotent in progress", http.StatusConflict)
				return
			}
			http.Error(w, "idempotent lock error", http.StatusInternalServerError)
			return
		}

		defer func(store Store, ctx context.Context, key string) {
			err := store.Unlock(ctx, key)
			if err != nil {
				logger.Error("failed to unlock idempotent key",
					zap.Error(err), zap.String("key", key))
			}
		}(store, r.Context(), key)

		cached, err := store.Get(r.Context(), key)
		if err != nil {
			http.Error(w, "failed to retrieve cache response", http.StatusInternalServerError)
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
		err = store.Set(ctx, key, &Response{
			Status: recorderWriter.status,
			Header: recorderWriter.cloneHeaders(),
			Body:   recorderWriter.body.Bytes(),
		}, 5*time.Minute)
		if err != nil {
			logger.Error("failed to store response",
				zap.Error(err),
				zap.String("key", key),
			)
		}
	})
}

func serveCachedResponse(cached *Response, w http.ResponseWriter) {
	copyAndReplaceHeaders(w.Header(), cached.Header)
	cStatus := cached.Status
	if cStatus > 0 {
		// Only if status is set write it, else default behaviour would result in 200
		w.WriteHeader(cStatus)
	}
	// Errors are likely due to client disconnect.
	_, _ = w.Write(cached.Body)
}

func copyAndReplaceHeaders(dst, src http.Header) {
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
