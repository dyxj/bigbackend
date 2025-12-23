package userprofile_test

import (
	"context"

	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/entity"
	"github.com/dyxj/bigbackend/pkg/sqldb"
	"github.com/stretchr/testify/mock"
)

type mockCreatorRepo struct {
	mock.Mock
}

func (m *mockCreatorRepo) InsertUserProfile(
	ctx context.Context, tx sqldb.Executable, input entity.UserProfile,
) (entity.UserProfile, error) {
	returnArgs := m.Called(ctx, tx, input)
	return returnArgs.Get(0).(entity.UserProfile), returnArgs.Error(1)
}
