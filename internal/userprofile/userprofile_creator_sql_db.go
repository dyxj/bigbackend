package userprofile

import (
	"context"

	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/entity"
	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/table"
	"github.com/dyxj/bigbackend/pkg/audit"
	"github.com/dyxj/bigbackend/pkg/sqldb"
	"github.com/go-jet/jet/v2/postgres"
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

	inputAuditable := userProfileAuditableEntity{P: &input}
	audit.InitInsertFields(inputAuditable)

	stmt := c.buildStatement(input)

	_, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return entity.UserProfile{}, err
	}

	return input, nil
}

func (c *CreatorSQLDB) buildStatement(input entity.UserProfile) postgres.InsertStatement {
	return table.UserProfile.
		INSERT(table.UserProfile.AllColumns).
		MODEL(input)
}
