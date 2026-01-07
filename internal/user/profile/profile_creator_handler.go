package profile

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dyxj/bigbackend/pkg/errorx"
	"github.com/dyxj/bigbackend/pkg/httpx"
	"github.com/dyxj/bigbackend/pkg/sqldb"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Creator interface {
	CreateUserProfileTx(ctx context.Context, tx sqldb.Executable, input UserProfile) (UserProfile, error)
}

type CreatorHandler struct {
	logger  *zap.Logger
	tm      sqldb.TransactionManager
	creator Creator
	mapper  Mapper
}

func NewCreatorHandler(
	logger *zap.Logger,
	tm sqldb.TransactionManager,
	creator Creator,
	mapper Mapper,
) *CreatorHandler {
	return &CreatorHandler{logger: logger, tm: tm, creator: creator, mapper: mapper}
}

func (c *CreatorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var cRequest CreateRequest
	err := json.NewDecoder(r.Body).Decode(&cRequest)
	if err != nil {
		c.logger.Warn("failed to decode create user profile request", zap.Error(err))
		httpx.BadRequestResponse("invalid request body",
			map[string]string{"error": err.Error()},
			w)
		return
	}

	vErr := cRequest.Validate()
	if vErr != nil {
		c.logger.Warn("create user profile request validation failed", zap.Any("userId", cRequest.UserID))
		httpx.ValidationFailedResponse(vErr, w)
		return
	}

	userId := chi.URLParam(r, "id")
	if userId != cRequest.UserID.String() {
		c.logger.Warn("user ID in URL does not match user ID in request body",
			zap.String("urlUserId", userId),
			zap.String("bodyUserId", cRequest.UserID.String()))
		httpx.BadRequestResponse("user ID in URL does not match user ID in request body", nil, w)
		return
	}

	tx, err := c.tm.BeginTx(r.Context(), nil)
	if err != nil {
		c.logger.Error("failed to begin transaction", zap.Error(err))
		httpx.InternalServerErrorResponse("", w)
		return
	}
	defer func() {
		err := tx.Rollback()
		if err != nil && !errors.Is(err, sql.ErrTxDone) {
			c.logger.Error("failed to rollback transaction", zap.Error(err))
		}
	}()

	input := c.mapper.CreateRequestToModel(cRequest)

	created, err := c.creator.CreateUserProfileTx(r.Context(), tx, input)
	if err != nil {
		c.resolveError(err, w)
		return
	}

	err = tx.Commit()
	if err != nil {
		c.logger.Error("failed to commit transaction", zap.Error(err))
		httpx.InternalServerErrorResponse("", w)
		return
	}

	response := c.mapper.ModelToResponse(created)

	httpx.JsonResponse(http.StatusCreated, response, w)
}

func (c *CreatorHandler) resolveError(err error, w http.ResponseWriter) {
	var uErr *errorx.UniqueViolationError
	if errors.As(err, &uErr) {
		c.logger.Warn("failed to insert user profile due to unique key violation", zap.Error(uErr))
		httpx.ConflictResponse("user profile already exists", nil, w)
		return
	}
	var vErr *errorx.ValidationError
	if errors.As(err, &vErr) {
		c.logger.Warn("failed to insert user profile due to validation error", zap.Error(vErr))
		httpx.ValidationFailedResponse(vErr, w)
		return

	}
	c.logger.Error("failed to insert user profile due to internal server error", zap.Error(err))
	httpx.InternalServerErrorResponse("", w)
}
