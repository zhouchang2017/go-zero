package sqlx

import (
	"context"
	"sync"
)

type (
	SqlHook interface {
		OnExecStart(ctx context.Context, q string, args ...interface{})
		OnExecFinished(ctx context.Context, err error)

		OnExecStmtStart(ctx context.Context, q string, args ...interface{})
		OnExecStmtFinished(ctx context.Context, err error)

		OnQueryStart(ctx context.Context, q string, args ...interface{})
		OnQueryFinished(ctx context.Context, err error)

		OnQueryStmtStart(ctx context.Context, q string, args ...interface{})
		OnQueryStmtFinished(ctx context.Context, err error)
	}
)

var sqlHooks []SqlHook
var _lock sync.Mutex

func RegisterSqlHooks(h SqlHook) {
	_lock.Lock()
	defer _lock.Unlock()
	sqlHooks = append(sqlHooks, h)
}

func rangeSqlHookOnStart(command string, ctx context.Context, q string, args ...interface{}) {
	for _, hook := range sqlHooks {
		if hook != nil {
			switch command {
			case "exec":
				hook.OnExecStart(ctx, q, args...)
			case "execStmt":
				hook.OnExecStmtStart(ctx, q, args...)
			case "query":
				hook.OnQueryStart(ctx, q, args...)
			case "queryStmt":
				hook.OnQueryStmtStart(ctx, q, args...)
			}
		}
	}
}

func rangeSqlHookOnFinished(command string, ctx context.Context, err error) {
	for _, hook := range sqlHooks {
		if hook != nil {
			switch command {
			case "exec":
				hook.OnExecFinished(ctx, err)
			case "execStmt":
				hook.OnExecStmtFinished(ctx, err)
			case "query":
				hook.OnQueryFinished(ctx, err)
			case "queryStmt":
				hook.OnQueryStmtFinished(ctx, err)
			}
		}
	}
}
