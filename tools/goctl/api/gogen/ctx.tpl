package internal

import "context"

type ctxKey int

const (
	actionKey ctxKey = iota
	requestIdKey
)

//CtxGetAction 从context中获取action
func CtxGetAction(ctx context.Context) (string, bool) {
	action, ok := ctx.Value(actionKey).(string)
	return action, ok
}

//CtxWithAction 将action写入context进行传递
func CtxWithAction(ctx context.Context, action string) context.Context {
	return context.WithValue(ctx, actionKey, action)
}

//CtxGetRequestId 从context中获取requestId
func CtxGetRequestId(ctx context.Context) (string, bool) {
	requestId, ok := ctx.Value(requestIdKey).(string)
	return requestId, ok
}

//CtxWithRequestId 将requestId写入context进行传递
func CtxWithRequestId(ctx context.Context, requestId string) context.Context {
	return context.WithValue(ctx, requestIdKey, requestId)
}