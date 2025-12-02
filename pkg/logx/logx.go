package logx

import (
	"go.uber.org/zap"
)

func InitLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()

	config.OutputPaths = []string{"stdout"}
	config.InitialFields = map[string]interface{}{
		"service": "bigbackend",
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}
