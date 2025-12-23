package sqldb

import "github.com/lib/pq"

// IsUniqueViolationError TODO consider extracting violation key
func IsUniqueViolationError(err *pq.Error) bool {
	if err.Code.Name() == "unique_violation" {
		return true
	}
	return false
}
