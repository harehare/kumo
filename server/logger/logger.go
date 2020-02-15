package logger

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

func NewLogger() *zap.Logger {
	// TODO:
	logger, _ := zap.NewDevelopment()

	return logger
}
