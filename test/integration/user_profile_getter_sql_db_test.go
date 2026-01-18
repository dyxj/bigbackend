//go:build integration

package integration

import (
	"testing"
	"time"

	"github.com/dyxj/bigbackend/internal/user/profile"
	"github.com/dyxj/bigbackend/pkg/errorx"
	"github.com/dyxj/bigbackend/pkg/logx"
	"github.com/dyxj/bigbackend/pkg/testx"
	"github.com/dyxj/bigbackend/test"
	"github.com/dyxj/bigbackend/test/faker"
	"github.com/stretchr/testify/assert"
)

func TestSQL_GetUserProfileByUserID(t *testing.T) {
	logger, err := logx.InitLogger()
	if err != nil {
		t.Fatalf("failed to initialize logger: %v", err)
	}
	dbConn := testx.GlobalEnv().DBConn()

	t.Run("should get user profile successfully", func(t *testing.T) {
		ctx := t.Context()

		t.Cleanup(func() {
			test.TruncateUserProfile(dbConn)
		})

		creator := profile.NewCreatorSQLDB(logger)
		getter := profile.NewGetterSQLDB(logger, dbConn)

		uProfile := faker.UserProfileEntity()

		inserted, err := creator.InsertUserProfile(ctx, dbConn, uProfile)
		if err != nil {
			t.Fatalf("failed to insert user profile: %v", err)
		}

		found, err := getter.FindUserProfileByUserID(ctx, uProfile.UserID)
		if err != nil {
			t.Fatalf("failed to get user profile by id: %v", err)
		}

		assert.Equal(t, inserted.ID, found.ID)
		assert.WithinDuration(t, inserted.CreateTime, found.CreateTime, time.Second)
		assert.WithinDuration(t, inserted.UpdateTime, found.UpdateTime, time.Second)
		assert.Equal(t, inserted.Version, found.Version)

		assert.Equal(t, uProfile.UserID, found.UserID)
		assert.Equal(t, uProfile.FirstName, found.FirstName)
		assert.Equal(t, uProfile.LastName, found.LastName)
		assert.Equal(t, uProfile.DateOfBirth, found.DateOfBirth)
	})

	t.Run("should return not found error", func(t *testing.T) {
		ctx := t.Context()

		t.Cleanup(func() {
			test.TruncateUserProfile(dbConn)
		})

		getter := profile.NewGetterSQLDB(logger, dbConn)

		_, err := getter.FindUserProfileByUserID(ctx, faker.UserProfile().UserID)
		if err == nil {
			t.Fatalf("expected error but got nil")
		}

		assert.ErrorIs(t, err, errorx.ErrNotFound)
	})
}
