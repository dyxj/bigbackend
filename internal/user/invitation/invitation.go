package invitation

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending  Status = "PENDING"
	StatusAccepted Status = "ACCEPTED"
	StatusExpired  Status = "EXPIRED"
)

var Statuses = []Status{
	StatusPending,
	StatusAccepted,
	StatusExpired,
}

type UserInvitation struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	StatusRaw  Status    `json:"statusRaw"`
	ExpiryTime time.Time `json:"expiryTime"`
	Token      string    `json:"token"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
	Version    int32     `json:"version"`
}

func (u *UserInvitation) Status() Status {
	if u.StatusRaw == StatusExpired || u.StatusRaw == StatusAccepted {
		return u.StatusRaw
	}

	if time.Now().After(u.ExpiryTime) {
		return StatusExpired
	}

	return StatusPending
}

func (u *UserInvitation) IsValidForCreate() bool {
	return u.IsValid()
}

func (u *UserInvitation) IsValid() bool {
	return u.Email != "" &&
		u.Token != ""
}
