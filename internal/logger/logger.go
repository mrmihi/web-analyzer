package logger

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Field represents a log field
type Field struct {
	Key   string
	Value interface{}
}

// Logger interface
type Logger interface {
	DebugCtx(ctx context.Context, msg string, fields ...Field)
	InfoCtx(ctx context.Context, msg string, fields ...Field)
	WarnCtx(ctx context.Context, msg string, fields ...Field)
	ErrorCtx(ctx context.Context, msg string, fields ...Field)
}

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

// StubLogger implements the Logger interface for testing
// Not used in test as of now in the project, but can be useful for unit tests
// TODO: Use StubLogger in unit tests
type StubLogger struct {
	sync.RWMutex
	lastMessage string
}

// NewStubLogger creates a new StubLogger
func NewStubLogger() *StubLogger {
	return &StubLogger{}
}

// DebugCtx logs a debug message with context
func (sl *StubLogger) DebugCtx(ctx context.Context, msg string, fields ...Field) {
	sl.log(ctx, "debug", msg, fields...)
}

// InfoCtx logs an info message with context
func (sl *StubLogger) InfoCtx(ctx context.Context, msg string, fields ...Field) {
	sl.log(ctx, "information", msg, fields...)
}

// WarnCtx logs a warning message with context
func (sl *StubLogger) WarnCtx(ctx context.Context, msg string, fields ...Field) {
	sl.log(ctx, "warning", msg, fields...)
}

// ErrorCtx logs an error message with context
func (sl *StubLogger) ErrorCtx(ctx context.Context, msg string, fields ...Field) {
	sl.log(ctx, "error", msg, fields...)
}

// LastMessage returns the last logged message
func (sl *StubLogger) LastMessage() string {
	sl.RLock()
	defer sl.RUnlock()

	return sl.lastMessage
}

// log is a helper function to log messages
func (sl *StubLogger) log(ctx context.Context, level, msg string, fields ...Field) {
	line := fmt.Sprintf("%s: %s", level, msg)

	sl.Lock()
	defer sl.Unlock()

	sl.lastMessage = line
}

// Global logger instance
var globalLogger Logger

// SetLogger sets the global logger instance
func SetLogger(logger Logger) {
	globalLogger = logger
}

// GetLogger returns the global logger instance
func GetLogger() Logger {
	if globalLogger == nil {
		// Default to development logger if not set
		globalLogger = NewZapLogger(true)
	}
	return globalLogger
}

// Helper functions to use the global logger

// DebugCtx logs a debug message with context
func DebugCtx(ctx context.Context, msg string, fields ...Field) {
	GetLogger().DebugCtx(ctx, msg, fields...)
}

// InfoCtx logs an info message with context
func InfoCtx(ctx context.Context, msg string, fields ...Field) {
	GetLogger().InfoCtx(ctx, msg, fields...)
}

// WarnCtx logs a warning message with context
func WarnCtx(ctx context.Context, msg string, fields ...Field) {
	GetLogger().WarnCtx(ctx, msg, fields...)
}

// ErrorCtx logs an error message with context
func ErrorCtx(ctx context.Context, msg string, fields ...Field) {
	GetLogger().ErrorCtx(ctx, msg, fields...)
}
