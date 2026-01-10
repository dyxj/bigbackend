package account

import (
	"time"

	"github.com/google/uuid"
)

type UserAccount struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreateTime   time.Time `json:"createTime"`
	UpdateTime   time.Time `json:"updateTime"`
	Version      int32     `json:"version"`
}
