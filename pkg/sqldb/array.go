package sqldb

import (
	"database/sql/driver"

	"github.com/lib/pq"
)

type Array[T any] []T

func (A *Array[T]) Value() (driver.Value, error) {
	return pq.Array(A).Value()
}

func (A *Array[T]) Scan(src any) error {
	return pq.Array(A).Scan(src)
}
