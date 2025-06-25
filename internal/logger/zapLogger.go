package logger

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLogger implements the Logger interface using Uber's Zap
type ZapLogger struct {
	logger *zap.Logger
}

// NewZapLogger creates a new ZapLogger with the given configuration
func NewZapLogger(development bool) *ZapLogger {
	var logger *zap.Logger
	var err error

	if development {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, err = config.Build()
	} else {
		config := zap.NewProductionConfig()
		logger, err = config.Build()
	}

	if err != nil {
		panic(err)
	}

	return &ZapLogger{
		logger: logger,
	}
}

// Convert custom fields to zap fields
func fieldsToZapFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Value)
	}
	return zapFields
}

// DebugCtx logs a debug message with context
func (l *ZapLogger) DebugCtx(ctx context.Context, msg string, fields ...Field) {
	l.logger.Debug(msg, fieldsToZapFields(fields)...)
}

// InfoCtx logs an info message with context
func (l *ZapLogger) InfoCtx(ctx context.Context, msg string, fields ...Field) {
	l.logger.Info(msg, fieldsToZapFields(fields)...)
}

// WarnCtx logs a warning message with context
func (l *ZapLogger) WarnCtx(ctx context.Context, msg string, fields ...Field) {
	l.logger.Warn(msg, fieldsToZapFields(fields)...)
}

// ErrorCtx logs an error message with context
func (l *ZapLogger) ErrorCtx(ctx context.Context, msg string, fields ...Field) {
	l.logger.Error(msg, fieldsToZapFields(fields)...)
}
