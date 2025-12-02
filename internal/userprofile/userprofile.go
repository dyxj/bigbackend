package userprofile

import (
	"time"

	"cloud.google.com/go/civil"
	"github.com/google/uuid"
)

type UserProfile struct {
	ID          uuid.UUID  `json:"id"`
	Email       string     `json:"email"`
	FirstName   string     `json:"firstName"`
	LastName    string     `json:"lastName"`
	DateOfBirth civil.Date `json:"dateOfBirth"`
	CreateTime  time.Time  `json:"createTime"`
	UpdateTime  time.Time  `json:"updateTime"`
	Version     int32      `json:"version"`
}
