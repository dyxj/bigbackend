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

type createConfig struct {
	defaultExpiryDuration time.Duration
	invitationURL         string
}

type CreateOption func(*createConfig)

func CreateOptDefaultExpiryDuration(d time.Duration) CreateOption {
	return func(c *createConfig) {
		c.defaultExpiryDuration = d
	}
}

func CreateOptInvitationURL(invitationURL string) CreateOption {
	return func(c *createConfig) {
		c.invitationURL = invitationURL
	}
}

const sysDefaultExpiryDuration = time.Hour * 24
const sysDefaultInvitationURL = "http://localhost:8080/invitation?token="

type creator struct {
	logger      *zap.Logger
	tm          sqldb.TransactionManager
	creatorRepo CreatorRepo
	mapper      Mapper
	publisher   EventPublisher
	expirer     Expirer
	cfg         createConfig
}

func NewCreator(
	logger *zap.Logger,
	tm sqldb.TransactionManager,
	creatorRepo CreatorRepo,
	mapper Mapper,
	publisher EventPublisher,
	expirer Expirer,
	option ...CreateOption,
) Creator {
	cfg := createConfig{
		defaultExpiryDuration: sysDefaultExpiryDuration,
		invitationURL:         sysDefaultInvitationURL,
	}

	for _, opt := range option {
		opt(&cfg)
	}

	return &creator{
		logger: logger, tm: tm, creatorRepo: creatorRepo, mapper: mapper,
		publisher: publisher, expirer: expirer, cfg: cfg,
	}
}

func (c *creator) CreateUserInvitation(
	ctx context.Context, input UserInvitation,
) (UserInvitation, error) {

	isValid := input.IsValidForCreate()
	if !isValid {
		return UserInvitation{}, &errorx.ValidationError{}
	}

	input.ExpiryTime = time.Now().Add(c.cfg.defaultExpiryDuration)
	input.StatusRaw = StatusPending
	input.Token = uuid.New().String()

	// TODO check if account is created for email

	tx, err := c.tm.BeginTx(ctx, nil)
	if err != nil {
		c.logger.Error("failed to begin transaction", zap.Error(err))
		return UserInvitation{}, err
	}
	defer sqldb.TxRollback(tx, c.logger)

	err = c.expirer.ExpireInvitationsByEmailTx(ctx, tx, input.Email)
	if err != nil {
		return UserInvitation{}, err
	}

	entityInput := c.mapper.ModelToEntity(input)
	createdEntity, err := c.creatorRepo.InsertUserInvitation(ctx, tx, entityInput)
	if err != nil {
		return UserInvitation{}, err
	}

	err = c.publisher.Publish(ctx, tx)
	if err != nil {
		return UserInvitation{}, err
	}

	err = tx.Commit()
	if err != nil {
		c.logger.Error("failed to commit transaction", zap.Error(err))
		return UserInvitation{}, err
	}

	return c.mapper.EntityToModel(createdEntity), nil
}

type Creator interface {
	CreateUserInvitation(ctx context.Context, input UserInvitation) (UserInvitation, error)
}

type CreatorRepo interface {
	InsertUserInvitation(
		ctx context.Context, tx sqldb.Executable, input entity.UserInvitation,
	) (entity.UserInvitation, error)
}

type EventPublisher interface {
	Publish(ctx context.Context, tx sqldb.Executable) error
}

type Expirer interface {
	ExpireInvitationsByEmailTx(ctx context.Context, tx sqldb.Executable, email string) error
}
