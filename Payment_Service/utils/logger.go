package utils

import (
	"context"
	"os"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func InitLogger() *zap.Logger {

	loggerCfg := zap.NewProductionConfig()

	loggerCfg.EncoderConfig.TimeKey = "ts"
	loggerCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	l, _ := loggerCfg.Build()
	logger = l
	return logger
}

func AppLogger() *zap.Logger {
	if logger == nil {
		return zap.NewNop()
	}
	return logger
}

func WithTraceContext(ctx context.Context) *zap.Logger {
	if logger == nil {
		InitLogger()
	}
	span := trace.SpanContextFromContext(ctx)
	if !span.IsValid() {
		return logger
	}
	return logger.With(
		zap.String("trace_id", span.TraceID().String()),
		zap.String("span_id", span.SpanID().String()),
	)
}

func ShutdownLogger() {
	if logger != nil {
		_ = logger.Sync()
	}
}

func SetLevelFromEnv() {
	lvl := os.Getenv("LOG_LEVEL")
	if lvl == "" || logger == nil {
		return
	}
	var level zapcore.Level
	if err := level.Set(lvl); err == nil {
		_ = zap.ReplaceGlobals(logger.WithOptions(zap.IncreaseLevel(level)))
	}
}
