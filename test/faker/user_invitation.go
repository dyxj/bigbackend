package faker

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/entity"
	"github.com/dyxj/bigbackend/internal/user/invitation"
	"github.com/google/uuid"
)

func UserInvitation() invitation.UserInvitation {
	return invitation.UserInvitation{
		ID:         uuid.New(),
		Email:      gofakeit.Email(),
		StatusRaw:  RandomType(invitation.Statuses),
		ExpiryTime: futureDate(),
		Token:      gofakeit.UUID(),
		CreateTime: gofakeit.Date(),
		UpdateTime: gofakeit.Date(),
		Version:    0,
	}
}

func UserInvitationEntity() entity.UserInvitation {
	return entity.UserInvitation{
		ID:         uuid.New(),
		Email:      gofakeit.Email(),
		Status:     string(RandomType(invitation.Statuses)),
		ExpiryTime: futureDate(),
		Token:      gofakeit.UUID(),
		CreateTime: gofakeit.Date(),
		UpdateTime: gofakeit.Date(),
		Version:    0,
	}
}
