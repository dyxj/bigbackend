package userprofile

import (
	"strings"
	"time"

	"cloud.google.com/go/civil"
	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/entity"
	"github.com/google/uuid"
)

type UserProfile struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"userId"`
	FirstName   string     `json:"firstName"`
	LastName    string     `json:"lastName"`
	DateOfBirth civil.Date `json:"dateOfBirth"`
	CreateTime  time.Time  `json:"createTime"`
	UpdateTime  time.Time  `json:"updateTime"`
	Version     int32      `json:"version"`
}

func (u *UserProfile) Sanitize() {
	u.FirstName = strings.Trim(u.FirstName, " ")
	u.LastName = strings.Trim(u.LastName, " ")
}

func (u *UserProfile) IsValidForCreate() bool {
	return u.isValid()
}

func (u *UserProfile) IsValidForUpdate() bool {
	return u.isValid() &&
		u.ID != uuid.Nil
}

func (u *UserProfile) isValid() bool {
	return u.FirstName != "" &&
		u.LastName != "" &&
		u.DateOfBirth.IsValid() &&
		!u.DateOfBirth.IsZero() &&
		u.DateOfBirth.Before(civil.DateOf(time.Now()))
}

// userProfileAuditableEntity adapts entity.UserProfile to repo.Auditable.
type userProfileAuditableEntity struct{ P *entity.UserProfile }

func (a userProfileAuditableEntity) GetID() uuid.UUID          { return a.P.ID }
func (a userProfileAuditableEntity) SetID(id uuid.UUID)        { a.P.ID = id }
func (a userProfileAuditableEntity) SetCreateTime(t time.Time) { a.P.CreateTime = t }
func (a userProfileAuditableEntity) SetUpdateTime(t time.Time) { a.P.UpdateTime = t }
func (a userProfileAuditableEntity) GetVersion() int32         { return a.P.Version }
func (a userProfileAuditableEntity) SetVersion(v int32)        { a.P.Version = v }
