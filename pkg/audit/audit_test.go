package audit

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type RandomEntity struct {
	ID         uuid.UUID `json:"id"`
	Message    string    `json:"message"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
	Version    int32     `json:"version"`
}

type RandomEntityAuditable struct {
	P *RandomEntity
}

func (a RandomEntityAuditable) GetID() uuid.UUID          { return a.P.ID }
func (a RandomEntityAuditable) SetID(id uuid.UUID)        { a.P.ID = id }
func (a RandomEntityAuditable) SetCreateTime(t time.Time) { a.P.CreateTime = t }
func (a RandomEntityAuditable) SetUpdateTime(t time.Time) { a.P.UpdateTime = t }
func (a RandomEntityAuditable) GetVersion() int32         { return a.P.Version }
func (a RandomEntityAuditable) SetVersion(v int32)        { a.P.Version = v }

func TestInitInsertFields(t *testing.T) {
	t.Run("should init all fields for new entity", func(t *testing.T) {
		entity := &RandomEntity{Message: "hello"}
		auditable := RandomEntityAuditable{P: entity}

		InitInsertFields(auditable)

		assert.NotEqual(t, uuid.Nil, auditable.P.ID)
		assert.WithinDuration(t, time.Now(), auditable.P.CreateTime, time.Second)
		assert.WithinDuration(t, time.Now(), auditable.P.UpdateTime, time.Second)
		assert.Equal(t, int32(1), auditable.P.Version)
	})

	t.Run("should overwrite existing ID", func(t *testing.T) {
		existingID := uuid.New()
		entity := &RandomEntity{ID: existingID, Message: "world"}
		auditable := RandomEntityAuditable{P: entity}

		InitInsertFields(auditable)

		assert.NotEqual(t, existingID, auditable.P.ID)
		assert.WithinDuration(t, time.Now(), auditable.P.CreateTime, time.Second)
		assert.WithinDuration(t, time.Now(), auditable.P.UpdateTime, time.Second)
		assert.Equal(t, int32(1), auditable.P.Version)
	})
}

func TestInitUpdateFields(t *testing.T) {
	t.Run("should update fields for existing entity", func(t *testing.T) {
		entity := &RandomEntity{
			ID:         uuid.New(),
			Message:    "update test",
			CreateTime: time.Now().Add(-time.Hour),
			UpdateTime: time.Now().Add(-time.Hour),
			Version:    2,
		}
		auditable := RandomEntityAuditable{P: entity}

		InitUpdateFields(auditable)

		assert.WithinDuration(t, time.Now(), auditable.P.UpdateTime, time.Second)
		assert.Equal(t, int32(3), auditable.P.Version)
	})
}
