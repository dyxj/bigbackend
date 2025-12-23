package faker

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

func pastDate() time.Time {
	duration := time.Duration(24 * gofakeit.IntRange(1, 50))
	return time.Now().Add(time.Hour * -duration)
}
