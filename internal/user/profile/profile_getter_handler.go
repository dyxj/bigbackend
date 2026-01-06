package profile

import (
	"context"
	"errors"
	"net/http"

	"github.com/dyxj/bigbackend/pkg/errorx"
	"github.com/dyxj/bigbackend/pkg/httpx"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type GetterHandler struct {
	logger *zap.Logger
	getter Getter
	mapper Mapper
}

func NewGetterHandler(logger *zap.Logger, getter Getter, mapper Mapper) *GetterHandler {
	return &GetterHandler{logger: logger, getter: getter, mapper: mapper}
}

func (g *GetterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		g.logger.Warn("failed to parse id", zap.String("id", idStr), zap.Error(err))
		httpx.BadRequestResponse("invalid id",
			map[string]string{"error": err.Error()},
			w)
		return
	}

	profile, err := g.getter.GetUserProfileByUserID(r.Context(), id)
	if err != nil {
		g.resolveError(err, w)
		return
	}

	resp := g.mapper.ModelToResponse(profile)

	httpx.JsonResponse(http.StatusOK, resp, w)
}

func (g *GetterHandler) resolveError(err error, w http.ResponseWriter) {
	if errors.Is(err, errorx.ErrNotFound) {
		httpx.NotFoundResponse(w)
		return
	}
	httpx.InternalServerErrorResponse("", w)
}

type Getter interface {
	GetUserProfileByUserID(ctx context.Context, userID uuid.UUID) (UserProfile, error)
}
