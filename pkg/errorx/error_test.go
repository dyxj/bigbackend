package errorx

import (
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

		expectedMessage := "validation error: field1:must not be empty field2:must be a valid email"
		assert.Equal(t, expectedMessage, err.Error())
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
				"field1": "abc",
				"field2": "123",
			},
		}

		expectedMessage := "unique violation error: field1:abc field2:123"
		assert.Equal(t, expectedMessage, err.Error())
	})

	t.Run("without properties", func(t *testing.T) {
		err := &UniqueViolationError{}

		expectedMessage := "unique violation error"
		assert.Equal(t, expectedMessage, err.Error())
	})
}
