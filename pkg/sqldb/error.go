package sqldb

import "github.com/lib/pq"

func IsUniqueViolationError(err *pq.Error) bool {
	if err.Code.Name() == "unique_violation" {
		return true
	}
	return false
}
