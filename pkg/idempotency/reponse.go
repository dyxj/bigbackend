package idempotency

import (
	"bytes"
	"net/http"
	"time"
)

// Response a cache version of an HTTP response.ß
type Response struct {
	Status   int
	Header   http.Header
	Body     []byte
	StoredAt time.Time
}

type responseRecorderWriter struct {
	w                    http.ResponseWriter
	status               int
	body                 *bytes.Buffer
	hasWroteStatusHeader bool
}

func newResponseRecorderWriter(w http.ResponseWriter) *responseRecorderWriter {
	return &responseRecorderWriter{w: w, body: &bytes.Buffer{}}
}

func (r *responseRecorderWriter) Header() http.Header {
	return r.w.Header()
}

func (r *responseRecorderWriter) Write(b []byte) (int, error) {
	// bytes.Buffer.Write does not return an error, rather panics if buffer is too large
	_, _ = r.body.Write(b)
	return r.w.Write(b)
}

func (r *responseRecorderWriter) WriteHeader(statusCode int) {
	// Ensure only a single write for status, as per net/http implementation
	if !r.hasWroteStatusHeader {
		r.hasWroteStatusHeader = true
		r.status = statusCode
	}
	r.w.WriteHeader(statusCode)
}

func (r *responseRecorderWriter) cloneHeaders() http.Header {
	ori := r.w.Header()
	clone := make(http.Header, len(ori))
	for k, vv := range ori {
		copyV := make([]string, len(vv))
		copy(copyV, vv)
		clone[k] = copyV
	}
	return clone
}
