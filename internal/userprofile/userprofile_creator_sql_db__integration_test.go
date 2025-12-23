//go:build integration

package userprofile_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/dyxj/bigbackend/internal/faker"
	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/entity"
	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/table"
	"github.com/dyxj/bigbackend/internal/userprofile"
	"github.com/dyxj/bigbackend/pkg/logx"
	"github.com/dyxj/bigbackend/pkg/sqldb"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreatorSQLDB_InsertUserProfile(t *testing.T) {
	getwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	fmt.Println(getwd)
	logger, err := logx.InitLogger()
	if err != nil {
		t.Fatalf("failed to initialize logger: %v", err)
	}
	dbConn := sqldb.SetupTestDB(t)

	t.Run("should insert successfully", func(t *testing.T) {
		t.Cleanup(func() {
			truncateUserProfile(dbConn)
		})

		creator := userprofile.NewCreatorSQLDB(logger)

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
			truncateUserProfile(dbConn)
		})

		creator := userprofile.NewCreatorSQLDB(logger)

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

		var pqErr *pq.Error
		isExpectedError := assert.ErrorAs(t, err, &pqErr, "expected pq.Error type")
		if isExpectedError {
			assert.Equal(t, "23505", string(pqErr.Code), "expected unique_violation error code")
			assert.Equal(t, fmt.Sprintf("Key (user_id)=(%s) already exists.", input2.UserID), pqErr.Detail, "expected unique violation detail")
		}
	})
}

func truncateUserProfile(dbConn *sql.DB) {
	log.Printf("truncating user_profile table")
	_, err := dbConn.Exec("TRUNCATE TABLE user_profile CASCADE;")
	if err != nil {
		log.Printf("failed to truncate user_profile table: %v", err)
		return
	}
}
