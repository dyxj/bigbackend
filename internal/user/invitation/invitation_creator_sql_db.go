package invitation

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

// InsertUserInvitation inserts a new user invitation into the database.
func (c *CreatorSQLDB) InsertUserInvitation(
	ctx context.Context,
	tx sqldb.Executable,
	input entity.UserInvitation,
) (entity.UserInvitation, error) {
	c.logger.Debug("inserting user invitation", zap.Any("email", input.Email))

	inputAuditable := userInvitationAuditableEntity{E: &input}
	audit.SetInsertFields(inputAuditable)

	stmt := c.buildStatement(input)

	_, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return entity.UserInvitation{}, c.resolveError(err)
	}

	c.logger.Debug("inserted user invitation", zap.Any("email", input.Email))

	return input, nil
}

func (c *CreatorSQLDB) buildStatement(input entity.UserInvitation) postgres.InsertStatement {
	return table.UserInvitation.
		INSERT(table.UserInvitation.AllColumns).
		MODEL(input)
}

func (c *CreatorSQLDB) resolveError(err error) error {
	var pqErr *pq.Error
	isPqErr := errors.As(err, &pqErr)
	if isPqErr && sqldb.IsUniqueViolationError(pqErr) {
		if pqErr.Constraint == dbcUkToken {
			return &errorx.UniqueViolationError{
				Properties: map[string]string{
					"token": "token already exists",
				},
			}
		}
		if pqErr.Constraint == dbcUkAcceptedPendingEmail {
			return &errorx.UniqueViolationError{
				Properties: map[string]string{
					"email": "email already has a pending or accepted invitation",
				},
			}
		}
		return &errorx.UniqueViolationError{}
	}
	return err
}
