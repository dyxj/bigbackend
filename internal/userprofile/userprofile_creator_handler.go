package userprofile

import (
	"net/http"

	"github.com/dyxj/bigbackend/pkg/sqldb"
	"go.uber.org/zap"
)

type CreatorHandler struct {
	logger  *zap.Logger
	tm      sqldb.TransactionManager
	creator *Creator
}

func NewCreatorHandler(logger *zap.Logger, tm sqldb.TransactionManager, creator *Creator) *CreatorHandler {
	return &CreatorHandler{logger: logger, tm: tm, creator: creator}
}

func (c *CreatorHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

}
