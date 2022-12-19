package logx

import (
	"context"
)

type ctxKey int

const requestIdKey = "RequestId"
const (
	loggerKey ctxKey = iota
	reqIdKey
)

// GetRequestId 获取上下文的RequestId
func GetRequestId(ctx context.Context) string {
	id, ok := ctx.Value(reqIdKey).(string)
	if ok {
		return id
	}
	return ""
}

// SetRequestId 设置RequestId
func SetRequestId(ctx context.Context, l Logger, reqId string) (context.Context, Logger) {
	logger := l.WithRequestId(reqId)
	return WithCtx(context.WithValue(ctx, reqIdKey, reqId), logger), logger
}

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
