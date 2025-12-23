package sqldb

import (
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestIsUniqueViolationError(t *testing.T) {
	t.Run("should return true for unique violation error", func(t *testing.T) {
		err := &pq.Error{
			Code: "23505", // unique_violation
		}
		assert.True(t, IsUniqueViolationError(err))
	})

	t.Run("should return false for other errors", func(t *testing.T) {
		err := &pq.Error{
			Code: "42P01", // undefined_table
		}
		assert.False(t, IsUniqueViolationError(err))
	})
}
