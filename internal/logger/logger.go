package logger

import (
	"context"
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
