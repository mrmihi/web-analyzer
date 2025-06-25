package logger

import "context"

var globalLogger Logger

// SetLogger sets the global logger instance
func SetLogger(logger Logger) {
	globalLogger = logger
}

// GetLogger returns the global logger instance
func GetLogger() Logger {
	if globalLogger == nil {
		globalLogger = NewZapLogger(true)
	}
	return globalLogger
}

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
