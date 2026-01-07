package invitation

import (
	"context"
	"errors"
	"time"

	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/entity"
	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/table"
	"github.com/dyxj/bigbackend/pkg/errorx"
	"github.com/dyxj/bigbackend/pkg/sqldb"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
)

type UpdaterSQLDB struct {
	logger *zap.Logger
}

func NewUpdaterSQLDB(logger *zap.Logger) *UpdaterSQLDB {
	return &UpdaterSQLDB{
		logger: logger,
	}
}

func (u *UpdaterSQLDB) UpdateInvitationTx(
	ctx context.Context, tx sqldb.Queryable, input entity.UserInvitation,
) (entity.UserInvitation, error) {

	stmt := table.UserInvitation.
		UPDATE(
			table.UserInvitation.AllColumns.
				Except(table.UserInvitation.CreateTime),
		).
		MODEL(input).
		WHERE(postgres.AND(
			table.UserInvitation.ID.EQ(postgres.UUID(input.ID)),
			table.UserInvitation.Version.EQ(postgres.Int32(input.Version)),
		)).
		RETURNING(table.UserInvitation.AllColumns)

	var updated entity.UserInvitation
	err := stmt.QueryContext(ctx, tx, &updated)
	if err != nil {
		return entity.UserInvitation{}, u.resolveUpdateInvitationTxErr(err)
	}

	return updated, nil
}

func (u *UpdaterSQLDB) resolveUpdateInvitationTxErr(err error) error {
	if errors.Is(err, qrm.ErrNoRows) {
		return errorx.ErrNotFound
	}
	return err
}

func (u *UpdaterSQLDB) BatchUpdateInvitationTx(
	ctx context.Context, tx sqldb.Queryable, inputs []entity.UserInvitation,
) ([]entity.UserInvitation, error) {

	var rows []postgres.RowExpression
	for _, input := range inputs {
		rows = append(rows, postgres.WRAP(
			postgres.UUID(input.ID),
			postgres.String(input.Email),
			postgres.String(input.Status),
			postgres.TimestampzT(input.ExpiryTime),
			postgres.String(input.Token),
			postgres.Int32(input.Version),
		))
	}

	uId := postgres.StringColumn("id")
	uEmail := postgres.StringColumn("email")
	uStatus := postgres.StringColumn("status")
	uExpiry := postgres.TimestampzColumn("expiry")
	uToken := postgres.StringColumn("token")
	uOldVersion := postgres.IntegerColumn("old_version")

	updateData := postgres.VALUES(rows...).
		AS("u", uId, uEmail, uStatus, uExpiry, uToken, uOldVersion)

	now := time.Now()

	stmt := table.UserInvitation.
		UPDATE().
		SET(
			table.UserInvitation.Email.SET(uEmail),
			table.UserInvitation.Status.SET(uStatus),
			table.UserInvitation.ExpiryTime.SET(uExpiry),
			table.UserInvitation.Version.SET(
				table.UserInvitation.Version.ADD(postgres.Int32(1)),
			),
			table.UserInvitation.Token.SET(uToken),
			table.UserInvitation.UpdateTime.SET(postgres.TimestampzT(now)),
		).
		FROM(updateData).
		WHERE(
			postgres.AND(
				table.UserInvitation.ID.EQ(uId),
				table.UserInvitation.Version.EQ(uOldVersion),
			),
		).
		RETURNING(table.UserInvitation.AllColumns)

	var result []entity.UserInvitation
	err := stmt.QueryContext(ctx, tx, &result)
	if err != nil {
		return nil, err
	}

	if len(result) != len(inputs) {
		return result, errorx.ErrConflict
	}

	return result, nil
}
