package rescue

import (
	"github.com/zeromicro/go-zero/core/logx"
	"runtime/debug"
)

// Recover is used with defer to do cleanup on panics.
// Use it like:
//
//	defer Recover(func() {})
func Recover(cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		logx.GlobalLogger().Errorf("%v\n%s", p, string(debug.Stack()))
	}
}
