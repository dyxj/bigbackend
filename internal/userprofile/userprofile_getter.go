package userprofile

import (
	"context"

	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/entity"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type getter struct {
	logger     *zap.Logger
	getterRepo GetterRepo
	mapper     Mapper
}

func NewGetter(logger *zap.Logger, getterRepo GetterRepo, mapper Mapper) Getter {
	return &getter{logger: logger, getterRepo: getterRepo, mapper: mapper}
}

func (g *getter) GetUserProfileByUserID(ctx context.Context, userID uuid.UUID) (UserProfile, error) {
	profile, err := g.getterRepo.FindUserProfileByUserID(ctx, userID)

	if err != nil {
		return UserProfile{}, err
	}

	return g.mapper.EntityToModel(profile), nil
}

type GetterRepo interface {
	FindUserProfileByUserID(ctx context.Context, userID uuid.UUID) (entity.UserProfile, error)
}
