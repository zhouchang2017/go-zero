package logx

import (
	"context"
)

type ctxKey int

const requestIdKey = "RequestId"
const (
	loggerKey ctxKey = iota
)

// WithCtx 添加到上下文
func WithCtx(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

// FromCtx 从上下文中获取默认日志
func FromCtx(ctx context.Context) Logger {
	if ctx == nil {
		return globalLogger
	}
	logger, ok := ctx.Value(loggerKey).(Logger)
	if ok {
		return logger
	}
	return globalLogger
}
