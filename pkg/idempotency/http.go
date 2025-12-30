package idempotency

import (
	"net/http"
	"time"
)

const DefaultKey = "Idempotency-Key"

type KeyExtractor func(r *http.Request) string

func DefaultKeyExtractor(r *http.Request) string {
	return r.Header.Get(DefaultKey)
}

// Middleware returns an http.Handler that enforces idempotency using the provided store.
//
// Requests without a key are passed through.
//
// If a request contains an idempotency key, the middleware will:
//   - Serve a cached response if present.
//
// Else check if idempotency in progress
//
//   - if yes, wait for "retryWaitDuration" then try to fetch from store and replace response.
//     if store does not contain cache response, return downstream response without storing.
//
//   - if no, record downstream handler response, store it, and return it.
func Middleware(
	extractor KeyExtractor,
	store Store,
	retryWaitDuration time.Duration,
	next http.Handler,
) http.Handler {
	if extractor == nil {
		extractor = DefaultKeyExtractor
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := extractor(r)
		if key == "" {
			next.ServeHTTP(w, r)
			return
		}

		if cached, ok := store.Get(key); ok {
			serveCachedResponse(cached, w)
			return
		}

		// run next with response copier
		recorderWriter := newResponseRecorderWriter(w)
		next.ServeHTTP(recorderWriter, r)

		isInProgress, ok := r.Context().Value(isInProgressKey).(bool)
		if ok && isInProgress {

			<-time.After(retryWaitDuration)

			cached, ok := store.Get(key)
			if key != "" && ok {
				serveCachedResponse(cached, w)
			}

			return
		}

		// cache response
		store.Set(key,
			&Response{
				Status:   recorderWriter.status,
				Header:   recorderWriter.cloneHeaders(),
				Body:     recorderWriter.body.Bytes(),
				StoredAt: time.Now(),
			})
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
