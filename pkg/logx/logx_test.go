package logx

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitLogger(t *testing.T) {
	// Keep original stdout and restore it when the test finishes
	old := os.Stdout
	defer func() { os.Stdout = old }()

	// Create a pipe to capture output
	r, w, err := os.Pipe()
	require.NoError(t, err, "os.Pipe should not return an error")
	os.Stdout = w

	// Initialize logger - this will now write to our pipe
	logger, err := InitLogger()
	require.NoError(t, err, "InitLogger should not return an error")
	require.NotNil(t, logger, "InitLogger should return a non-nil logger")

	// Log a message
	logger.Info("test message")

	// Close the writer
	err = w.Close()
	require.NoError(t, err, "Closing the pipe writer should not return an error")

	// Read the output from the pipe
	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	require.NoError(t, err)

	// Assert that the output contains the service name and the message
	output := buf.String()
	assert.Contains(t, output, `"service":"bigbackend"`, "Log output should contain the service tag")
	assert.Contains(t, output, `"msg":"test message"`, "Log output should contain the message")
}
