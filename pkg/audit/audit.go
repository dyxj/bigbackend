package audit

import (
	"time"

	"github.com/google/uuid"
)

type Auditable interface {
	GetID() uuid.UUID
	SetID(uuid.UUID)
	SetCreateTime(time.Time)
	SetUpdateTime(time.Time)
	GetVersion() int32
	SetVersion(int32)
}

func InitInsertFields[T Auditable](t T) {
	now := time.Now()

	t.SetID(uuid.New())
	t.SetCreateTime(now)
	t.SetUpdateTime(now)
	t.SetVersion(1)
}

func InitUpdateFields[T Auditable](t T) {
	t.SetUpdateTime(time.Now())
	t.SetVersion(t.GetVersion() + 1)
}
