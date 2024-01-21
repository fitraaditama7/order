package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

var logger *zap.Logger
var once sync.Once

func Log() *zap.Logger {
	once.Do(func() {
		config := zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.DisableStacktrace = true

		logger, _ = config.Build()
	})
	return logger
}
