package userprofile

import (
	"context"
	"errors"

	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/entity"
	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/table"
	"github.com/dyxj/bigbackend/pkg/audit"
	"github.com/dyxj/bigbackend/pkg/errorx"
	"github.com/dyxj/bigbackend/pkg/sqldb"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

type CreatorSQLDB struct {
	logger *zap.Logger
}

func NewCreatorSQLDB(logger *zap.Logger) *CreatorSQLDB {
	return &CreatorSQLDB{
		logger: logger,
	}
}

// InsertUserProfile inserts a new user profile into the database.
// Ignores and automatically sets ID, CreateTime, UpdateTime, Version fields from input.
func (c *CreatorSQLDB) InsertUserProfile(
	ctx context.Context,
	tx sqldb.Executable,
	input entity.UserProfile,
) (entity.UserProfile, error) {
	c.logger.Debug("inserting user profile", zap.Any("userId", input.UserID))

	inputAuditable := userProfileAuditableEntity{P: &input}
	audit.InitInsertFields(inputAuditable)

	stmt := c.buildStatement(input)

	_, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return entity.UserProfile{}, c.resolveError(err, input)
	}

	c.logger.Debug("inserted user profile", zap.Any("userId", input.UserID))

	return input, nil
}

func (c *CreatorSQLDB) buildStatement(input entity.UserProfile) postgres.InsertStatement {
	return table.UserProfile.
		INSERT(table.UserProfile.AllColumns).
		MODEL(input)
}

func (c *CreatorSQLDB) resolveError(err error, input entity.UserProfile) error {
	var pqErr *pq.Error
	isPqErr := errors.As(err, &pqErr)
	if isPqErr && sqldb.IsUniqueViolationError(pqErr) {
		c.logger.Warn("failed to insert user profile due to unique key violation",
			zap.Any("userId", input.UserID),
			zap.String("detail", pqErr.Detail),
		)
		return &errorx.UniqueViolationError{
			Properties: map[string]string{"userId": input.UserID.String()},
		}
	}
	return err
}
