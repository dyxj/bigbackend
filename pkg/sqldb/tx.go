package sqldb

import (
	"database/sql"
	"errors"

	"go.uber.org/zap"
)

func TxRollback(tx *sql.Tx, logger *zap.Logger) {
	err := tx.Rollback()
	if err != nil && !errors.Is(err, sql.ErrTxDone) {
		logger.WithOptions(zap.AddCallerSkip(1)).
			Error("failed to rollback transaction", zap.Error(err))
	}
}
