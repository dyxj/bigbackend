package profile

import (
	"time"

	"cloud.google.com/go/civil"
	"github.com/dyxj/bigbackend/pkg/errorx"
	"github.com/google/uuid"
)

type CreateRequest struct {
	UserID      uuid.UUID  `json:"userId"`
	FirstName   string     `json:"firstName"`
	LastName    string     `json:"lastName"`
	DateOfBirth civil.Date `json:"dateOfBirth"`
}

func (r *CreateRequest) Validate() *errorx.ValidationError {
	errors := make(map[string]string)

	if r.UserID == uuid.Nil {
		errors["userId"] = "is required"
	}
	if r.FirstName == "" {
		errors["firstName"] = "is required"
	}
	if r.LastName == "" {
		errors["lastName"] = "is required"
	}
	if !r.DateOfBirth.IsValid() || r.DateOfBirth.IsZero() || !r.DateOfBirth.Before(civil.DateOf(time.Now())) {
		errors["dateOfBirth"] = "is invalid or in the future"
	}

	if len(errors) > 0 {
		return &errorx.ValidationError{Properties: errors}
	}

	return nil
}

type UpdateRequest struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"userId"`
	FirstName   string     `json:"firstName"`
	LastName    string     `json:"lastName"`
	DateOfBirth civil.Date `json:"dateOfBirth"`
	CreateTime  time.Time  `json:"createTime"`
	UpdateTime  time.Time  `json:"updateTime"`
	Version     int32      `json:"version"`
}

type Response struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"userId"`
	FirstName   string     `json:"firstName"`
	LastName    string     `json:"lastName"`
	DateOfBirth civil.Date `json:"dateOfBirth"`
	CreateTime  time.Time  `json:"createTime"`
	UpdateTime  time.Time  `json:"updateTime"`
	Version     int32      `json:"version"`
}
