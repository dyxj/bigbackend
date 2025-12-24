package userprofile

import (
	"context"

	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/entity"
	"github.com/dyxj/bigbackend/pkg/errorx"
	"github.com/dyxj/bigbackend/pkg/sqldb"
	"go.uber.org/zap"
)

type creator struct {
	logger      *zap.Logger
	creatorRepo CreatorRepo
	mapper      Mapper
}

func NewCreator(logger *zap.Logger, creatorRepo CreatorRepo, mapper Mapper) Creator {
	return &creator{logger: logger, creatorRepo: creatorRepo, mapper: mapper}
}

func (c *creator) CreateUserProfileTx(ctx context.Context, tx sqldb.Executable, input UserProfile) (UserProfile, error) {

	input.Sanitize()

	isValid := input.IsValidForCreate()
	if !isValid {
		return UserProfile{}, &errorx.ValidationError{}
	}

	userProfileEntity := c.mapper.ModelToEntity(input)
	createdEntity, err := c.creatorRepo.InsertUserProfile(ctx, tx, userProfileEntity)
	if err != nil {
		return UserProfile{}, err
	}

	return c.mapper.EntityToModel(createdEntity), nil
}

type CreatorRepo interface {
	InsertUserProfile(ctx context.Context, tx sqldb.Executable, input entity.UserProfile) (entity.UserProfile, error)
}
