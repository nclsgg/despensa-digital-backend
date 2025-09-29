package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Setup inicializa o logger global do projeto.
func Setup(mode string) func() {
	var (
		cfg    zap.Config
		logger *zap.Logger
		err    error
	)

	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "release", "production":
		cfg = zap.NewProductionConfig()
	default:
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	logger, err = cfg.Build()
	if err != nil {
		panic("failed to initialize zap logger: " + err.Error())
	}

	undo := zap.ReplaceGlobals(logger)
	zap.RedirectStdLog(logger)

	return func() {
		_ = logger.Sync()
		undo()
	}
}
