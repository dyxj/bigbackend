package mapx

import (
	"time"

	"cloud.google.com/go/civil"
	"github.com/google/uuid"
)

func MapTime(t time.Time) time.Time {
	return t
}

func MapDate(date civil.Date) civil.Date {
	return date
}

func MapUUID(id uuid.UUID) uuid.UUID {
	return id
}
