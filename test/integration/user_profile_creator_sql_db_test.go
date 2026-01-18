//go:build integration

package integration

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/entity"
	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/table"
	"github.com/dyxj/bigbackend/internal/user/profile"
	"github.com/dyxj/bigbackend/pkg/errorx"
	"github.com/dyxj/bigbackend/pkg/logx"
	"github.com/dyxj/bigbackend/pkg/testx"
	"github.com/dyxj/bigbackend/test"
	"github.com/dyxj/bigbackend/test/faker"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/stretchr/testify/assert"
)

func TestCreatorSQLDB_InsertUserProfile(t *testing.T) {
	logger, err := logx.InitLogger()
	if err != nil {
		t.Fatalf("failed to initialize logger: %v", err)
	}
	dbConn := testx.GlobalEnv().DBConn()

	t.Run("should insert successfully", func(t *testing.T) {
		t.Cleanup(func() {
			test.TruncateUserProfile(dbConn)
		})

		creator := profile.NewCreatorSQLDB(logger)

		tx, err := dbConn.Begin()
		if err != nil {
			t.Fatalf("failed to begin transaction: %v", err)
		}
		defer func() {
			err := tx.Rollback()
			if err != nil && !errors.Is(err, sql.ErrTxDone) {
				t.Fatalf("failed to rollback transaction: %v", err)
			}
		}()

		input := faker.UserProfileEntity()
		if err != nil {
			t.Fatalf("failed to generate fake user profile: %v", err)
		}

		inserted, err := creator.InsertUserProfile(
			context.Background(),
			tx,
			input,
		)

		err = tx.Commit()
		if err != nil {
			t.Fatalf("failed to commit transaction: %v", err)
		}

		var selected entity.UserProfile
		err = table.UserProfile.
			SELECT(table.UserProfile.AllColumns).
			WHERE(table.UserProfile.ID.EQ(postgres.UUID(inserted.ID))).
			Query(dbConn, &selected)
		if err != nil {
			t.Fatalf("failed to select inserted user profile: %v", err)
		}

		assert.NotNil(t, inserted.ID, "inserted ID should not be nil")
		assert.WithinDuration(t, time.Now(), inserted.CreateTime, time.Second, "inserted CreateTime should be recent")
		assert.WithinDuration(t, time.Now(), inserted.UpdateTime, time.Second, "inserted UpdateTime should be recent")
		assert.Equal(t, int32(1), inserted.Version, "inserted Version should be 1")

		assert.Equal(t, inserted.ID, selected.ID)
		assert.WithinDuration(t, inserted.CreateTime, selected.CreateTime, time.Second)
		assert.WithinDuration(t, inserted.UpdateTime, selected.UpdateTime, time.Second)

		assert.Equal(t, input.UserID, selected.UserID)
		assert.Equal(t, input.FirstName, selected.FirstName)
		assert.Equal(t, input.LastName, selected.LastName)
		assert.Equal(t, input.DateOfBirth, selected.DateOfBirth)
	})

	t.Run("should fail to insert with duplicate userId", func(t *testing.T) {
		t.Cleanup(func() {
			test.TruncateUserProfile(dbConn)
		})

		creator := profile.NewCreatorSQLDB(logger)

		tx, err := dbConn.Begin()
		if err != nil {
			t.Fatalf("failed to begin transaction: %v", err)
		}
		defer func() {
			err := tx.Rollback()
			if err != nil && !errors.Is(err, sql.ErrTxDone) {
				t.Fatalf("failed to rollback transaction: %v", err)
			}
		}()

		input := faker.UserProfileEntity()

		_, err = creator.InsertUserProfile(
			context.Background(),
			tx,
			input,
		)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		input2 := faker.UserProfileEntity()
		input2.UserID = input.UserID // duplicate userId

		_, err = creator.InsertUserProfile(
			context.Background(),
			tx,
			input2,
		)

		var uErr *errorx.UniqueViolationError
		isExpectedError := assert.ErrorAs(t, err, &uErr, "expected errorx.UniqueViolationError")
		if isExpectedError {
			assert.Equal(t, fmt.Sprintf("unique violation error | userId:%s", input2.UserID), uErr.Error(), "expected unique violation error message")
		}
	})
}
