package idempotency

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dyxj/bigbackend/test/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestHandler_NoKey(t *testing.T) {
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)

	mockStore := &MockIdempotencyStore{}
	mockErrWriter := &MockErrResponseWriter{}

	middleware := NewMiddleware(logger, mockStore,
		WithErrorResponseWriter(mockErrWriter.WriteError),
	)

	mockHandler := &faker.MockHandler{}
	mockHandler.On("ServeHTTP", mock.Anything, mock.Anything).
		Return()

	request := httptest.NewRequest(http.MethodDelete, "/", nil)
	respWriter := httptest.NewRecorder()

	middleware.Handler(mockHandler).
		ServeHTTP(respWriter, request)

	mockHandler.AssertNumberOfCalls(t, "ServeHTTP", 1)
	mockStore.AssertNumberOfCalls(t, "Lock", 0)
	mockStore.AssertNumberOfCalls(t, "Unlock", 0)
	mockStore.AssertNumberOfCalls(t, "Get", 0)
	mockStore.AssertNumberOfCalls(t, "Set", 0)
	mockErrWriter.AssertNumberOfCalls(t, "WriteError", 0)
}

func TestHandler_HTTP_Get(t *testing.T) {
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)

	mockStore := &MockIdempotencyStore{}
	mockErrWriter := &MockErrResponseWriter{}

	middleware := NewMiddleware(logger, mockStore,
		WithErrorResponseWriter(mockErrWriter.WriteError),
	)

	mockHandler := &faker.MockHandler{}
	mockHandler.On("ServeHTTP", mock.Anything, mock.Anything).
		Return()

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	respWriter := httptest.NewRecorder()
	request.Header.Set("Idempotency-Key", "fake-key")

	middleware.Handler(mockHandler).
		ServeHTTP(respWriter, request)

	mockHandler.AssertNumberOfCalls(t, "ServeHTTP", 1)
	mockStore.AssertNumberOfCalls(t, "Lock", 0)
	mockStore.AssertNumberOfCalls(t, "Unlock", 0)
	mockStore.AssertNumberOfCalls(t, "Get", 0)
	mockStore.AssertNumberOfCalls(t, "Set", 0)
	mockErrWriter.AssertNumberOfCalls(t, "WriteError", 0)
}

func TestHandler_WithKey_Run_Set(t *testing.T) {
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)
	key := "fake-key"

	mockStore := &MockIdempotencyStore{}
	mockErrWriter := &MockErrResponseWriter{}

	middleware := NewMiddleware(logger, mockStore,
		WithErrorResponseWriter(mockErrWriter.WriteError),
	)

	mockHandler := &faker.MockHandler{}
	mockHandler.On("ServeHTTP", mock.Anything, mock.Anything).
		Return()

	// Lock obtained
	mockStore.On("Lock",
		mock.Anything, key, mock.Anything).
		Return(nil)
	// Unlocked at end of middleware
	mockStore.On("Unlock", mock.Anything, key).
		Return(nil)
	// No cached response
	mockStore.On("Get", mock.Anything, key).
		Return(nil, nil)
	// Store the response
	mockStore.On("Set",
		mock.Anything, key, mock.Anything, mock.Anything).
		Return(nil)

	request := httptest.NewRequest(http.MethodPost, "/", nil)
	respWriter := httptest.NewRecorder()
	request.Header.Set("Idempotency-Key", key)

	middleware.Handler(mockHandler).
		ServeHTTP(respWriter, request)

	mockHandler.AssertNumberOfCalls(t, "ServeHTTP", 1)
	mockStore.AssertNumberOfCalls(t, "Lock", 1)
	mockStore.AssertNumberOfCalls(t, "Unlock", 1)
	mockStore.AssertNumberOfCalls(t, "Get", 1)
	mockStore.AssertNumberOfCalls(t, "Set", 1)
	mockErrWriter.AssertNumberOfCalls(t, "WriteError", 0)
}

func TestHandler_WithKey_Run_Set_UnlockFailed(t *testing.T) {
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)
	key := "fake-key"

	mockStore := &MockIdempotencyStore{}
	mockErrWriter := &MockErrResponseWriter{}

	middleware := NewMiddleware(logger, mockStore,
		WithErrorResponseWriter(mockErrWriter.WriteError),
	)

	mockHandler := &faker.MockHandler{}
	mockHandler.On("ServeHTTP", mock.Anything, mock.Anything).
		Return()

	// Lock obtained
	mockStore.On("Lock",
		mock.Anything, key, mock.Anything).
		Return(nil)
	// Unlocked at end of middleware
	mockStore.On("Unlock", mock.Anything, key).
		Return(errors.New("fake error"))
	// No cached response
	mockStore.On("Get", mock.Anything, key).
		Return(nil, nil)
	// Store the response
	mockStore.On("Set",
		mock.Anything, key, mock.Anything, mock.Anything).
		Return(nil)

	request := httptest.NewRequest(http.MethodPost, "/", nil)
	respWriter := httptest.NewRecorder()
	request.Header.Set("Idempotency-Key", key)

	middleware.Handler(mockHandler).
		ServeHTTP(respWriter, request)

	mockHandler.AssertNumberOfCalls(t, "ServeHTTP", 1)
	mockStore.AssertNumberOfCalls(t, "Lock", 1)
	mockStore.AssertNumberOfCalls(t, "Unlock", 1)
	mockStore.AssertNumberOfCalls(t, "Get", 1)
	mockStore.AssertNumberOfCalls(t, "Set", 1)
	mockErrWriter.AssertNumberOfCalls(t, "WriteError", 0)
}

func TestHandler_WithKey_Run_SetFailed(t *testing.T) {
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)
	key := "fake-key"

	mockStore := &MockIdempotencyStore{}
	mockErrWriter := &MockErrResponseWriter{}

	middleware := NewMiddleware(logger, mockStore,
		WithErrorResponseWriter(mockErrWriter.WriteError),
	)

	mockHandler := &faker.MockHandler{}
	mockHandler.On("ServeHTTP", mock.Anything, mock.Anything).
		Return()

	// Lock obtained
	mockStore.On("Lock",
		mock.Anything, key, mock.Anything).
		Return(nil)
	// Unlocked at end of middleware
	mockStore.On("Unlock", mock.Anything, key).
		Return(nil)
	// No cached response
	mockStore.On("Get", mock.Anything, key).
		Return(nil, nil)
	// Store the response
	mockStore.On("Set",
		mock.Anything, key, mock.Anything, mock.Anything).
		Return(errors.New("fake error"))

	request := httptest.NewRequest(http.MethodPost, "/", nil)
	respWriter := httptest.NewRecorder()
	request.Header.Set("Idempotency-Key", key)

	middleware.Handler(mockHandler).
		ServeHTTP(respWriter, request)

	mockHandler.AssertNumberOfCalls(t, "ServeHTTP", 1)
	mockStore.AssertNumberOfCalls(t, "Lock", 1)
	mockStore.AssertNumberOfCalls(t, "Unlock", 1)
	mockStore.AssertNumberOfCalls(t, "Get", 1)
	mockStore.AssertNumberOfCalls(t, "Set", 1)
	mockErrWriter.AssertNumberOfCalls(t, "WriteError", 0)
}

func TestHandler_WithKey_LockFailed(t *testing.T) {
	tt := []struct {
		name string
		err  error
	}{
		{
			name: "in progress",
			err:  ErrInProgress,
		},
		{
			name: "unknown error",
			err:  errors.New("fake error"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			logger, err := zap.NewDevelopment()
			assert.NoError(t, err)
			key := "fake-key"

			mockStore := &MockIdempotencyStore{}
			mockErrWriter := &MockErrResponseWriter{}

			middleware := NewMiddleware(logger, mockStore,
				WithErrorResponseWriter(mockErrWriter.WriteError),
			)

			mockHandler := &faker.MockHandler{}
			mockHandler.On("ServeHTTP", mock.Anything, mock.Anything).
				Return()

			// Lock obtained
			mockStore.On("Lock",
				mock.Anything, key, mock.Anything).
				Return(tc.err)
			// Expect error response to be written
			mockErrWriter.On("WriteError",
				mock.Anything, mock.Anything).
				Return()

			request := httptest.NewRequest(http.MethodPost, "/", nil)
			respWriter := httptest.NewRecorder()
			request.Header.Set("Idempotency-Key", key)

			middleware.Handler(mockHandler).
				ServeHTTP(respWriter, request)

			mockHandler.AssertNumberOfCalls(t, "ServeHTTP", 0)
			mockStore.AssertNumberOfCalls(t, "Lock", 1)
			mockStore.AssertNumberOfCalls(t, "Unlock", 0)
			mockStore.AssertNumberOfCalls(t, "Get", 0)
			mockStore.AssertNumberOfCalls(t, "Set", 0)
			mockErrWriter.AssertNumberOfCalls(t, "WriteError", 1)
		})
	}
}

func TestHandler_WithKey_CacheFound(t *testing.T) {
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)
	key := "fake-key"

	mockStore := &MockIdempotencyStore{}
	mockErrWriter := &MockErrResponseWriter{}

	middleware := NewMiddleware(logger, mockStore,
		WithErrorResponseWriter(mockErrWriter.WriteError),
	)

	mockHandler := &faker.MockHandler{}
	mockHandler.On("ServeHTTP", mock.Anything, mock.Anything).
		Return()

	// Lock obtained
	mockStore.On("Lock",
		mock.Anything, key, mock.Anything).
		Return(nil)
	// Unlocked at end of middleware
	mockStore.On("Unlock", mock.Anything, key).
		Return(nil)
	// No cached response
	mockStore.On("Get", mock.Anything, key).
		Return(&Response{
			Status: http.StatusCreated,
			Header: http.Header{"Content-Type": []string{"application/json"}},
			Body:   []byte(`{"message":"cached response"}`),
		}, nil)

	request := httptest.NewRequest(http.MethodPost, "/", nil)
	respWriter := httptest.NewRecorder()
	request.Header.Set("Idempotency-Key", key)

	respWriter.Header().Set("X-Fake-Header", "to-be-overwritten")
	middleware.Handler(mockHandler).
		ServeHTTP(respWriter, request)

	result := respWriter.Result()
	defer result.Body.Close()
	buf := &bytes.Buffer{}
	_, err = buf.ReadFrom(result.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, result.StatusCode)
	assert.Equal(t, "", result.Header.Get("X-Fake-Header"))
	assert.Equal(t, "application/json", result.Header.Get("Content-Type"))
	assert.Equal(t, `{"message":"cached response"}`, buf.String())

	mockHandler.AssertNumberOfCalls(t, "ServeHTTP", 0)
	mockStore.AssertNumberOfCalls(t, "Lock", 1)
	mockStore.AssertNumberOfCalls(t, "Unlock", 1)
	mockStore.AssertNumberOfCalls(t, "Get", 1)
	mockStore.AssertNumberOfCalls(t, "Set", 0)
	mockErrWriter.AssertNumberOfCalls(t, "WriteError", 0)
}

func TestHandler_WithKey_GetFailed(t *testing.T) {
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)
	key := "fake-key"

	mockStore := &MockIdempotencyStore{}
	mockErrWriter := &MockErrResponseWriter{}

	middleware := NewMiddleware(logger, mockStore,
		WithErrorResponseWriter(mockErrWriter.WriteError),
	)

	mockHandler := &faker.MockHandler{}
	mockHandler.On("ServeHTTP", mock.Anything, mock.Anything).
		Return()

	// Lock obtained
	mockStore.On("Lock",
		mock.Anything, key, mock.Anything).
		Return(nil)
	// Unlocked at end of middleware
	mockStore.On("Unlock", mock.Anything, key).
		Return(nil)
	// No cached response
	mockStore.On("Get", mock.Anything, key).
		Return(nil, errors.New("fake error"))
	// Expect error response to be written
	mockErrWriter.On("WriteError",
		mock.Anything, mock.Anything).
		Return()

	request := httptest.NewRequest(http.MethodPost, "/", nil)
	respWriter := httptest.NewRecorder()
	request.Header.Set("Idempotency-Key", key)

	middleware.Handler(mockHandler).
		ServeHTTP(respWriter, request)

	mockHandler.AssertNumberOfCalls(t, "ServeHTTP", 0)
	mockStore.AssertNumberOfCalls(t, "Lock", 1)
	mockStore.AssertNumberOfCalls(t, "Unlock", 1)
	mockStore.AssertNumberOfCalls(t, "Get", 1)
	mockStore.AssertNumberOfCalls(t, "Set", 0)
	mockErrWriter.AssertNumberOfCalls(t, "WriteError", 1)
}

func TestDefaultErrorResponseWriter(t *testing.T) {
	tt := []struct {
		name       string
		err        error
		errMessage string
		statusCode int
	}{
		{
			name:       "ErrInProgress",
			err:        ErrInProgress,
			errMessage: "request with the same idempotency key is already in progress\n",
			statusCode: http.StatusConflict,
		},
		{
			name:       "other error",
			err:        errors.New("unknown error"),
			errMessage: "internal server error\n",
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			DefaultErrorResponseWriter(tc.err, recorder)

			result := recorder.Result()
			defer result.Body.Close()
			buf := &bytes.Buffer{}
			_, err := buf.ReadFrom(result.Body)
			assert.NoError(t, err)

			assert.Equal(t, tc.statusCode, result.StatusCode)
			assert.Equal(t, tc.errMessage, buf.String())
		})
	}
}

func TestDefaultKeyExtractor(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/", nil)
	request.Header.Set("Idempotency-Key", "test-key")

	extractedKey := DefaultKeyExtractor(request)
	assert.Equal(t, "test-key", extractedKey)
}

func TestTenantAndIdempotencyKeyExtractor_WithKey(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/", nil)
	request.Header.Set("Idempotency-Key", "test-key")
	request.Header.Set("X-Tenant-ID", "tenant-123")

	extractedKey := TenantAndIdempotency(request)
	assert.Equal(t, "tenant-123::test-key", extractedKey)
}

func TestTenantAndIdempotencyKeyExtractor_NoKey(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/", nil)
	request.Header.Set("Idempotency-Key", "")
	request.Header.Set("X-Tenant-ID", "tenant-123")

	extractedKey := TenantAndIdempotency(request)
	assert.Equal(t, "", extractedKey)
}

type MockErrResponseWriter struct {
	mock.Mock
}

func (m *MockErrResponseWriter) WriteError(err error, w http.ResponseWriter) {
	m.Called(err, w)
}

type MockIdempotencyStore struct {
	mock.Mock
}

func (m *MockIdempotencyStore) Lock(ctx context.Context, key string, options ...LockOption) error {
	rArgs := m.Called(ctx, key, options)
	return rArgs.Error(0)
}

func (m *MockIdempotencyStore) Unlock(ctx context.Context, key string) error {
	rArgs := m.Called(ctx, key)
	return rArgs.Error(0)
}

func (m *MockIdempotencyStore) Get(ctx context.Context, key string) (*Response, error) {
	rArgs := m.Called(ctx, key)
	if rArgs.Get(0) == nil {
		return nil, rArgs.Error(1)
	}
	return rArgs.Get(0).(*Response), rArgs.Error(1)
}

func (m *MockIdempotencyStore) Set(
	ctx context.Context, key string, resp *Response, expiryDuration time.Duration,
) error {
	rArgs := m.Called(ctx, key, resp, expiryDuration)
	return rArgs.Error(0)
}
