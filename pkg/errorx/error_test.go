package errorx

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationError(t *testing.T) {
	t.Run("with properties", func(t *testing.T) {
		err := &ValidationError{
			Properties: map[string]string{
				"field1": "must not be empty",
				"field2": "must be a valid email",
			},
		}

		errMessages := strings.Split(err.Error(), " | ")
		assert.Contains(t, errMessages, "validation error")
		assert.Contains(t, errMessages, "field1:must not be empty")
		assert.Contains(t, errMessages, "field2:must be a valid email")
	})

	t.Run("without properties", func(t *testing.T) {
		err := &ValidationError{}

		expectedMessage := "validation error"
		assert.Equal(t, expectedMessage, err.Error())
	})
}

func TestUniqueViolationError(t *testing.T) {
	t.Run("with properties", func(t *testing.T) {
		err := &UniqueViolationError{
			Properties: map[string]string{
				"field1": "123 already exists",
				"field2": "abc already exists",
			},
		}

		errMessages := strings.Split(err.Error(), " | ")
		assert.Contains(t, errMessages, "unique violation error")
		assert.Contains(t, errMessages, "field1:123 already exists")
		assert.Contains(t, errMessages, "field2:abc already exists")
	})

	t.Run("without properties", func(t *testing.T) {
		err := &UniqueViolationError{}

		expectedMessage := "unique violation error"
		assert.Equal(t, expectedMessage, err.Error())
	})
}
