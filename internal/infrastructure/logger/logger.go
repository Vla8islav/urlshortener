package logger

import (
	"go.uber.org/zap"
)

func NewSugaredLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	return logger.Sugar()
}
