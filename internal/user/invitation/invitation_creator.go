package invitation

import (
	"context"
	"time"

	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/entity"
	"github.com/dyxj/bigbackend/pkg/errorx"
	"github.com/dyxj/bigbackend/pkg/sqldb"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const sysDefaultExpiryDuration = time.Hour * 24

type creator struct {
	logger                *zap.Logger
	creatorRepo           CreatorRepo
	mapper                Mapper
	defaultExpiryDuration time.Duration
}

func NewCreator(
	logger *zap.Logger,
	creatorRepo CreatorRepo,
	mapper Mapper,
	defaultExpiryDuration time.Duration,
) Creator {
	if defaultExpiryDuration <= 0 {
		defaultExpiryDuration = sysDefaultExpiryDuration
	}
	return &creator{
		logger: logger, creatorRepo: creatorRepo, mapper: mapper, defaultExpiryDuration: defaultExpiryDuration,
	}
}

func (c *creator) CreateUserInvitationTx(
	ctx context.Context, tx sqldb.Executable, input UserInvitation,
) (UserInvitation, error) {

	isValid := input.IsValidForCreate()
	if !isValid {
		return UserInvitation{}, &errorx.ValidationError{}
	}

	input.ExpiryTime = time.Now().Add(c.defaultExpiryDuration)
	input.StatusRaw = StatusPending
	input.Token = uuid.New().String()

	entityInput := c.mapper.ModelToEntity(input)
	createdEntity, err := c.creatorRepo.InsertUserInvitation(ctx, tx, entityInput)
	if err != nil {
		return UserInvitation{}, err
	}

	return c.mapper.EntityToModel(createdEntity), nil
}

type Creator interface {
	CreateUserInvitationTx(ctx context.Context, tx sqldb.Executable, input UserInvitation) (UserInvitation, error)
}

type CreatorRepo interface {
	InsertUserInvitation(
		ctx context.Context, tx sqldb.Executable, input entity.UserInvitation,
	) (entity.UserInvitation, error)
}
