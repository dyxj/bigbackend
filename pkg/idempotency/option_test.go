package idempotency

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithKeyExtractor(t *testing.T) {
	config := &Config{}
	assert.Nil(t, config.extractor)

	keyExtractor := func(r *http.Request) string { return "fake-extractor" }

	WithKeyExtractor(keyExtractor)(config)

	assert.NotNil(t, config.extractor)
}

func TestWithErrorResponseWriter(t *testing.T) {
	config := &Config{}
	assert.Nil(t, config.errRespWriter)

	errRespWriter := func(err error, w http.ResponseWriter) {}

	WithErrorResponseWriter(errRespWriter)(config)

	assert.NotNil(t, config.errRespWriter)
}

func TestWithCacheExpiry(t *testing.T) {
	config := &Config{}
	assert.Equal(t, time.Duration(0), config.cacheExpiry)

	WithCacheExpiry(2 * time.Hour)(config)

	assert.Equal(t, 2*time.Hour, config.cacheExpiry)
}
