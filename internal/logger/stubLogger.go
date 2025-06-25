package logger

import (
	"context"
	"fmt"
	"sync"
)

// TODO: Use StubLogger in unit tests, implement context to enable tracing

// StubLogger implements the Logger interface for testing
// Not used in test as of now in the project, but can be useful for unit tests
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
