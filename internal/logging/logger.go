package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger() {
	config := zap.NewDevelopmentConfig()

	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	config.OutputPaths = []string{"stdout"}
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

	logger, err := config.Build()

	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	Logger = logger
}
