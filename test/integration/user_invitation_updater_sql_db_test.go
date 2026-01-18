//go:build integration

package integration

import (
	"testing"
	"time"

	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/entity"
	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/table"
	"github.com/dyxj/bigbackend/internal/user/invitation"
	"github.com/dyxj/bigbackend/pkg/logx"
	"github.com/dyxj/bigbackend/pkg/sqldb"
	"github.com/dyxj/bigbackend/pkg/testx"
	"github.com/dyxj/bigbackend/test"
	"github.com/dyxj/bigbackend/test/faker"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/stretchr/testify/assert"
)

func TestUpdaterSQLDBUserInvitation_BatchUpdateInvitationTx(t *testing.T) {
	logger, err := logx.InitLogger()
	if err != nil {
		t.Fatalf("failed to initialize logger: %v", err)
	}
	dbConn := testx.GlobalEnv().DBConn()

	t.Run("should batch update successfully", func(t *testing.T) {
		ctx := t.Context()
		t.Cleanup(func() {
			test.TruncateUserInvitation(dbConn)
		})

		creator := invitation.NewCreatorSQLDB(logger)
		updater := invitation.NewUpdaterSQLDB(logger)

		i1 := faker.UserInvitationEntity()
		i1.Status = string(invitation.StatusPending)
		i2 := faker.UserInvitationEntity()
		i2.Status = string(invitation.StatusPending)

		inserted1, err := creator.InsertUserInvitation(ctx, dbConn, i1)
		assert.NoError(t, err)
		inserted2, err := creator.InsertUserInvitation(ctx, dbConn, i2)
		assert.NoError(t, err)

		update1 := faker.UserInvitationEntity()
		update1.ID = inserted1.ID
		update1.Version = inserted1.Version
		update1.Status = string(invitation.StatusAccepted)

		inserted2.Status = string(invitation.StatusExpired)

		result, err := updater.BatchUpdateInvitationTx(ctx, dbConn, []entity.UserInvitation{
			update1,
			inserted2,
		})
		assert.Equal(t, 2, len(result))
		assert.NoError(t, err)

		result1 := result[0]
		assert.Equal(t, inserted1.ID, result1.ID)
		assert.Equal(t, update1.Email, result1.Email)
		assert.Equal(t, string(invitation.StatusAccepted), result1.Status)
		assert.WithinDuration(t, update1.ExpiryTime, result1.ExpiryTime, time.Second)
		assert.Equal(t, update1.Token, result1.Token)
		assert.WithinDuration(t, inserted1.CreateTime, result1.CreateTime, time.Second)
		assert.Less(t, inserted1.UpdateTime, result1.UpdateTime)
		assert.Equal(t, int32(2), result1.Version)

		result2 := result[1]
		assert.Equal(t, inserted2.ID, result2.ID)
		assert.Equal(t, inserted2.Email, result2.Email)
		assert.Equal(t, string(invitation.StatusExpired), result2.Status)
		assert.WithinDuration(t, inserted2.ExpiryTime, result2.ExpiryTime, time.Second)
		assert.Equal(t, inserted2.Token, result2.Token)
		assert.WithinDuration(t, inserted2.CreateTime, result2.CreateTime, time.Second)
		assert.Less(t, inserted2.UpdateTime, result2.UpdateTime)
		assert.Equal(t, int32(2), result2.Version)
	})

	t.Run("should abort update if any failed", func(t *testing.T) {
		ctx := t.Context()
		t.Cleanup(func() {
			test.TruncateUserInvitation(dbConn)
		})

		creator := invitation.NewCreatorSQLDB(logger)
		updater := invitation.NewUpdaterSQLDB(logger)

		i1 := faker.UserInvitationEntity()
		i1.Status = string(invitation.StatusPending)
		i2 := faker.UserInvitationEntity()
		i2.Status = string(invitation.StatusPending)

		inserted1, err := creator.InsertUserInvitation(ctx, dbConn, i1)
		assert.NoError(t, err)

		update1 := faker.UserInvitationEntity()
		update1.ID = inserted1.ID
		update1.Status = string(invitation.StatusAccepted)
		update1.Version = inserted1.Version

		tx, err := dbConn.Begin()
		assert.NoError(t, err)
		defer sqldb.TxRollback(tx, logger)

		result, err := updater.BatchUpdateInvitationTx(ctx, tx, []entity.UserInvitation{
			update1,
			i2,
		})
		if err != nil {
			err := tx.Rollback()
			assert.NoError(t, err)
		} else {
			err := tx.Commit()
			assert.NoError(t, err)
			assert.FailNow(t, "expected conflict error")
		}

		assert.Equal(t, 1, len(result))

		var selected entity.UserInvitation
		err = table.UserInvitation.
			SELECT(table.UserInvitation.AllColumns).
			WHERE(table.UserInvitation.ID.EQ(postgres.UUID(inserted1.ID))).
			Query(dbConn, &selected)
		assert.NoError(t, err)

		assert.Equal(t, inserted1.ID, selected.ID)
		assert.Equal(t, int32(1), selected.Version)
		assert.Equal(t, string(invitation.StatusPending), selected.Status)

		assert.Equal(t, inserted1.Email, selected.Email)
		assert.WithinDuration(t, inserted1.ExpiryTime, selected.ExpiryTime, time.Second)
		assert.Equal(t, inserted1.Token, selected.Token)
		assert.WithinDuration(t, inserted1.CreateTime, selected.CreateTime, time.Second)
		assert.WithinDuration(t, inserted1.UpdateTime, selected.UpdateTime, time.Second)
	})
}
