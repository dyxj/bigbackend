package faker

import (
	"math/rand"
)

func RandomType[T any](val []T) T {
	return val[rand.Intn(len(val))]
}
