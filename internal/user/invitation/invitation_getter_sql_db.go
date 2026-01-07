package invitation

import (
	"context"

	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/entity"
	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/table"
	"github.com/dyxj/bigbackend/pkg/sqldb"
	"github.com/go-jet/jet/v2/postgres"
	"go.uber.org/zap"
)

type GetterSQLDB struct {
	logger *zap.Logger
	sqlQ   sqldb.Queryable
}

func NewGetterSQLDB(logger *zap.Logger, sqlQ sqldb.Queryable) *GetterSQLDB {
	return &GetterSQLDB{
		logger: logger,
		sqlQ:   sqlQ,
	}
}

func (g *GetterSQLDB) ListByEmailTx(ctx context.Context, tx sqldb.Queryable, email string) ([]entity.UserInvitation, error) {
	g.logger.Debug("selecting invitations by email", zap.String("email", email))

	stmt := table.UserInvitation.
		SELECT(table.UserInvitation.AllColumns).
		FROM(table.UserInvitation).
		WHERE(table.UserInvitation.Email.EQ(postgres.String(email))).
		ORDER_BY(table.UserInvitation.CreateTime.DESC())

	var results []entity.UserInvitation
	err := stmt.QueryContext(ctx, tx, &results)
	if err != nil {
		return nil, err
	}

	g.logger.Debug("selected invitations by email", zap.String("email", email))

	return results, nil
}
