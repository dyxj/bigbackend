package faker

import (
	"cloud.google.com/go/civil"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/entity"
	"github.com/google/uuid"
)

func UserProfileEntity() entity.UserProfile {
	return entity.UserProfile{
		ID:          uuid.New(),
		UserID:      uuid.New(),
		FirstName:   gofakeit.FirstName(),
		LastName:    gofakeit.LastName(),
		DateOfBirth: civil.DateOf(gofakeit.PastDate()),
		CreateTime:  gofakeit.Date(),
		UpdateTime:  gofakeit.Date(),
		Version:     0,
	}
}
