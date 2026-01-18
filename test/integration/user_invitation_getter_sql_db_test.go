//go:build integration

package integration

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/dyxj/bigbackend/internal/user/invitation"
	"github.com/dyxj/bigbackend/pkg/logx"
	"github.com/dyxj/bigbackend/pkg/testx"
	"github.com/dyxj/bigbackend/test"
	"github.com/dyxj/bigbackend/test/faker"
	"github.com/stretchr/testify/assert"
)

func TestGetterSQLDBUserInvitation_InsertUserInvitation(t *testing.T) {
	logger, err := logx.InitLogger()
	if err != nil {
		t.Fatalf("failed to initialize logger: %v", err)
	}
	dbConn := testx.GlobalEnv().DBConn()

	t.Run("should list invitations successfully", func(t *testing.T) {
		ctx := t.Context()
		t.Cleanup(func() {
			test.TruncateUserInvitation(dbConn)
		})

		getter := invitation.NewGetterSQLDB(logger, dbConn)
		creator := invitation.NewCreatorSQLDB(logger)

		email := gofakeit.Email()
		i1 := faker.UserInvitationEntity()
		i1.Email = email
		i1.Status = string(invitation.StatusPending)
		i2 := faker.UserInvitationEntity()
		i2.Email = email
		i2.Status = string(invitation.StatusExpired)

		inserted1, err := creator.InsertUserInvitation(ctx, dbConn, i1)
		assert.NoError(t, err)
		inserted2, err := creator.InsertUserInvitation(ctx, dbConn, i2)
		assert.NoError(t, err)

		result, err := getter.ListByEmailTx(ctx, dbConn, email)
		assert.NoError(t, err)

		assert.Equal(t, 2, len(result))
		assert.Equal(t, inserted2.ID, result[0].ID)
		assert.Equal(t, inserted1.ID, result[1].ID)
	})

	t.Run("should list no invitations", func(t *testing.T) {
		ctx := t.Context()
		t.Cleanup(func() {
			test.TruncateUserInvitation(dbConn)
		})

		getter := invitation.NewGetterSQLDB(logger, dbConn)
		creator := invitation.NewCreatorSQLDB(logger)

		email := "notfound@email.com"
		i1 := faker.UserInvitationEntity()
		i1.Status = string(invitation.StatusPending)
		i2 := faker.UserInvitationEntity()
		i2.Status = string(invitation.StatusExpired)

		_, err := creator.InsertUserInvitation(ctx, dbConn, i1)
		assert.NoError(t, err)
		_, err = creator.InsertUserInvitation(ctx, dbConn, i2)
		assert.NoError(t, err)

		result, err := getter.ListByEmailTx(ctx, dbConn, email)
		assert.NoError(t, err)

		assert.Equal(t, 0, len(result))
	})
}
