package logx

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder

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
