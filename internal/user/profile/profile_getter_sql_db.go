package profile

import (
	"context"
	"errors"

	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/entity"
	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/table"
	"github.com/dyxj/bigbackend/pkg/errorx"
	"github.com/dyxj/bigbackend/pkg/sqldb"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
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

// FindUserProfileByUserID retrieves a user profile by user ID from the database.
func (g *GetterSQLDB) FindUserProfileByUserID(
	ctx context.Context,
	userID uuid.UUID,
) (entity.UserProfile, error) {
	g.logger.Debug("get user profile", zap.Any("userId", userID))

	stmt := g.buildStatement(userID)

	var result entity.UserProfile
	err := stmt.QueryContext(ctx, g.sqlQ, &result)
	if err != nil {
		return entity.UserProfile{}, g.resolveError(err)
	}

	g.logger.Debug("found user profile", zap.Any("userId", userID))

	return result, nil
}

func (g *GetterSQLDB) buildStatement(userID uuid.UUID) postgres.SelectStatement {
	return table.UserProfile.
		SELECT(table.UserProfile.AllColumns).
		FROM(table.UserProfile).
		WHERE(table.UserProfile.UserID.EQ(postgres.UUID(userID)))
}

func (g *GetterSQLDB) resolveError(err error) error {
	if errors.Is(err, qrm.ErrNoRows) {
		return errorx.ErrNotFound
	}
	return err
}
