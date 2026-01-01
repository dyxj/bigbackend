package idempotency

import (
	"net/http"
	"testing"

	"github.com/dyxj/bigbackend/test/faker"
	"github.com/stretchr/testify/assert"
)

func TestResponseWriter_Header(t *testing.T) {
	mockWriter := &faker.MockHTTPResponseWriter{}
	mockWriter.On("Header").
		Return(http.Header{
			"Content-Type": []string{"application/json"},
			"X-Custom":     []string{"custom-value"},
		})

	writer := newResponseRecorderWriter(mockWriter)

	assert.Equal(t, "application/json", writer.Header().Get("Content-Type"))
	assert.Equal(t, "custom-value", writer.Header().Get("X-Custom"))
	mockWriter.AssertNumberOfCalls(t, "Header", 2)
}

func TestResponseWriter_Write(t *testing.T) {
	mockWriter := &faker.MockHTTPResponseWriter{}
	testData := []byte(`{"message":"hello world"}`)
	mockWriter.On("Write", testData).
		Return(len(testData), nil)

	writer := newResponseRecorderWriter(mockWriter)

	n, err := writer.Write(testData)

	assert.NoError(t, err)
	assert.Equal(t, len(testData), n)
	assert.Equal(t, testData, writer.body.Bytes())
	mockWriter.AssertNumberOfCalls(t, "Write", 1)
}

func TestResponseWriter_WriteHeader(t *testing.T) {
	mockWriter := &faker.MockHTTPResponseWriter{}
	statusCode := http.StatusCreated
	mockWriter.On("WriteHeader", statusCode).Return()
	mockWriter.On("WriteHeader", http.StatusOK).Return()

	writer := newResponseRecorderWriter(mockWriter)

	writer.WriteHeader(statusCode)
	writer.WriteHeader(http.StatusOK) // This should not overwrite the first status

	assert.Equal(t, statusCode, writer.status)
	mockWriter.AssertNumberOfCalls(t, "WriteHeader", 2)
}

func TestResponseWriter_CloneHeaders(t *testing.T) {
	originalHeaders := http.Header{
		"Content-Type": []string{"application/json"},
		"X-Custom":     []string{"custom-value"},
	}

	mockWriter := &faker.MockHTTPResponseWriter{}
	mockWriter.On("Header").
		Return(originalHeaders)

	writer := newResponseRecorderWriter(mockWriter)

	clonedHeaders := writer.cloneHeaders()

	assert.Equal(t, originalHeaders, clonedHeaders)

	// Modify cloned headers and ensure original is unaffected
	clonedHeaders.Set("X-Custom", "modified-value")

	assert.Equal(t, "custom-value", originalHeaders.Get("X-Custom"))
	assert.Equal(t, "modified-value", clonedHeaders.Get("X-Custom"))

	mockWriter.AssertNumberOfCalls(t, "Header", 1)
}
