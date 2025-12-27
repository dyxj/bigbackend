package faker

import (
	"context"

	"cloud.google.com/go/civil"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/entity"
	"github.com/dyxj/bigbackend/internal/userprofile"
	"github.com/dyxj/bigbackend/pkg/sqldb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

func UserProfile() userprofile.UserProfile {
	return userprofile.UserProfile{
		ID:          uuid.New(),
		UserID:      uuid.New(),
		FirstName:   gofakeit.FirstName(),
		LastName:    gofakeit.LastName(),
		DateOfBirth: civil.DateOf(pastDate()),
		CreateTime:  gofakeit.Date(),
		UpdateTime:  gofakeit.Date(),
		Version:     0,
	}
}

func UserProfileEntity() entity.UserProfile {
	return entity.UserProfile{
		ID:          uuid.New(),
		UserID:      uuid.New(),
		FirstName:   gofakeit.FirstName(),
		LastName:    gofakeit.LastName(),
		DateOfBirth: civil.DateOf(pastDate()),
		CreateTime:  gofakeit.Date(),
		UpdateTime:  gofakeit.Date(),
		Version:     0,
	}
}

func UserProfileCreateRequest() userprofile.CreateRequest {
	return userprofile.CreateRequest{
		UserID:      uuid.New(),
		FirstName:   gofakeit.FirstName(),
		LastName:    gofakeit.LastName(),
		DateOfBirth: civil.DateOf(pastDate()),
	}
}

type UserProfileCreatorRepoMock struct {
	mock.Mock
}

func (m *UserProfileCreatorRepoMock) InsertUserProfile(
	ctx context.Context, tx sqldb.Executable, input entity.UserProfile,
) (entity.UserProfile, error) {
	returnArgs := m.Called(ctx, tx, input)
	return returnArgs.Get(0).(entity.UserProfile), returnArgs.Error(1)
}

type UserProfileCreatorMock struct {
	mock.Mock
}

func (m *UserProfileCreatorMock) CreateUserProfileTx(
	ctx context.Context,
	tx sqldb.Executable,
	input userprofile.UserProfile,
) (userprofile.UserProfile, error) {
	returnArgs := m.Called(ctx, tx, input)
	return returnArgs.Get(0).(userprofile.UserProfile), returnArgs.Error(1)
}
