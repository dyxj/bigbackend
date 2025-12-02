package userprofile

import (
	"context"
	"time"

	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/entity"
	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/table"
	"github.com/dyxj/bigbackend/pkg/sqldb"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CreatorSQLDB struct {
	logger *zap.Logger
}

func NewCreatorRepo(logger *zap.Logger) *CreatorSQLDB {
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
	now := time.Now()

	input.ID = uuid.New()
	input.Version = 1
	input.CreateTime = now
	input.UpdateTime = now

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
