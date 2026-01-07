//go:build integration

package integration

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/entity"
	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/table"
	"github.com/dyxj/bigbackend/internal/user/invitation"
	"github.com/dyxj/bigbackend/pkg/logx"
	"github.com/dyxj/bigbackend/pkg/testx"
	"github.com/dyxj/bigbackend/test/faker"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/stretchr/testify/assert"
)

func TestCreatorSQLDBUserInvitation_InsertUserInvitation(t *testing.T) {
	logger, err := logx.InitLogger()
	assert.NoError(t, err)
	dbConn := testx.GlobalEnv().DBConn()

	t.Run("should insert successfully", func(t *testing.T) {
		t.Cleanup(func() {
			truncateUserInvitation(dbConn)
		})

		creator := invitation.NewCreatorSQLDB(logger)

		tx, err := dbConn.Begin()
		assert.NoError(t, err)
		defer func() {
			err := tx.Rollback()
			if err != nil && !errors.Is(err, sql.ErrTxDone) {
				t.Fatalf("failed to rollback transaction: %v", err)
			}
		}()

		input := faker.UserInvitationEntity()
		assert.NoError(t, err)

		inserted, err := creator.InsertUserInvitation(
			context.Background(),
			tx,
			input,
		)

		err = tx.Commit()
		assert.NoError(t, err)

		var selected entity.UserInvitation
		err = table.UserInvitation.
			SELECT(table.UserInvitation.AllColumns).
			WHERE(table.UserInvitation.ID.EQ(postgres.UUID(inserted.ID))).
			Query(dbConn, &selected)
		assert.NoError(t, err)

		assert.NotNil(t, inserted.ID, "inserted ID should not be nil")
		assert.WithinDuration(t, time.Now(), inserted.CreateTime, time.Second, "inserted CreateTime should be recent")
		assert.WithinDuration(t, time.Now(), inserted.UpdateTime, time.Second, "inserted UpdateTime should be recent")
		assert.Equal(t, int32(1), inserted.Version, "inserted Version should be 1")

		assert.Equal(t, inserted.ID, selected.ID)
		assert.WithinDuration(t, inserted.CreateTime, selected.CreateTime, time.Second)
		assert.WithinDuration(t, inserted.UpdateTime, selected.UpdateTime, time.Second)

		assert.Equal(t, input.Email, selected.Email)
		assert.Equal(t, input.Token, selected.Token)
		assert.Equal(t, input.Status, selected.Status)
		assert.WithinDuration(t, input.ExpiryTime, selected.ExpiryTime, time.Second)
	})
}
