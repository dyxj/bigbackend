package userprofile

import (
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type GetterHandler struct {
	logger *zap.Logger
}

func NewGetterHandler(logger *zap.Logger) *GetterHandler {
	return &GetterHandler{logger: logger}
}

func (g *GetterHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	idStr := request.PathValue("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		g.log().Info("invalid id", zap.String("id", idStr))
		http.Error(writer, "invalid id", http.StatusBadRequest)
		return
	}

	_, err = writer.Write([]byte(id.String()))
	if err != nil {
		g.log().Error("failed to write response", zap.Error(err))
		return
	}
	return
}

func (g *GetterHandler) log() *zap.Logger {
	return g.logger.With(zap.String("component", "userprofile_getter_handler"))
}
