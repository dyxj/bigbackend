package invitation

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dyxj/bigbackend/pkg/errorx"
	"github.com/dyxj/bigbackend/pkg/httpx"
	"go.uber.org/zap"
)

type CreatorHandler struct {
	logger  *zap.Logger
	creator Creator
	mapper  Mapper
}

func NewCreatorHandler(
	logger *zap.Logger,
	creator Creator,
	mapper Mapper,
) *CreatorHandler {
	return &CreatorHandler{logger: logger, creator: creator, mapper: mapper}
}

func (c *CreatorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()
	var cRequest CreateRequest
	err := json.NewDecoder(r.Body).Decode(&cRequest)
	if err != nil {
		c.logger.Warn("failed to decode create user invitation request", zap.Error(err))
		httpx.BadRequestResponse("invalid request body",
			map[string]string{"error": err.Error()},
			w)
		return
	}

	vErr := cRequest.Validate()
	if vErr != nil {
		c.logger.Warn("create user invitation request validation failed", zap.Any("email", cRequest.Email))
		httpx.ValidationFailedResponse(vErr, w)
		return
	}

	input := c.mapper.CreateRequestToModel(cRequest)
	_, err = c.creator.CreateUserInvitation(r.Context(), input)
	if err != nil {
		c.resolveError(err, w, cRequest)
		return
	}

	httpx.JsonResponse(http.StatusOK, CreateResponse{Email: cRequest.Email}, w)
}

func (c *CreatorHandler) resolveError(err error, w http.ResponseWriter, cr CreateRequest) {
	var uErr *errorx.UniqueViolationError
	if errors.As(err, &uErr) {
		c.logger.Warn("failed to create user invitation due to unique violation", zap.Error(uErr))
		httpx.JsonResponse(http.StatusOK, CreateResponse{Email: cr.Email}, w)
		return
	}
	var vErr *errorx.ValidationError
	if errors.As(err, &vErr) {
		c.logger.Warn("failed to create user invitation due to validation error", zap.Error(vErr))
		httpx.JsonResponse(http.StatusOK, CreateResponse{Email: cr.Email}, w)
		return
	}
	c.logger.Error("failed to insert user invitation", zap.Error(err))
	httpx.InternalServerErrorResponse("", w)
}
