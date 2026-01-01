package faker

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

type MockHTTPResponseWriter struct {
	mock.Mock
}

func (m *MockHTTPResponseWriter) Header() http.Header {
	rArgs := m.Called()
	return rArgs.Get(0).(http.Header)
}

func (m *MockHTTPResponseWriter) Write(bytes []byte) (int, error) {
	rArgs := m.Called(bytes)
	return rArgs.Int(0), rArgs.Error(1)
}

func (m *MockHTTPResponseWriter) WriteHeader(statusCode int) {
	m.Called(statusCode)
}
