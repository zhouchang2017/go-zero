package {{.PkgName}}
import (
	"{{.ModuleUrl}}/core/logx"
	"{{.ModuleUrl}}/rest/httpx"
	"net/http"
	"sync"
	{{.ImportPackages}}
)

var {{.Action}}Actions map[string]func(svcCtx *svc.ServiceContext) types.ActionHandler
var _{{.Action}}Lock sync.Mutex

// {{.RegisterActionName}} {{.Action}}路由注册入口
func {{.RegisterActionName}}(actionName string, handler func(svcCtx *svc.ServiceContext) types.ActionHandler) {
	_{{.Action}}Lock.Lock()
	defer _{{.Action}}Lock.Unlock()
	if {{.Action}}Actions == nil {
		{{.Action}}Actions = make(map[string]func(svcCtx *svc.ServiceContext) types.ActionHandler)
	}
	if handler == nil {
		logx.GlobalLogger().Errorf("{{.Action}} action: %s handler is nil.", actionName)
		return
	}

	if _, ok := {{.Action}}Actions[actionName]; ok {
		logx.GlobalLogger().Errorf("{{.Action}} action: %s already registed.", actionName)
	}

	{{.Action}}Actions[actionName] = handler
}

// {{.ActionHandlerName}} {{.Action}}http处理入口
func {{.ActionHandlerName}}(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		action, ok := internal.CtxGetAction(ctx)
		if !ok {
			httpx.ErrorCtx(ctx, w, types.ErrActionParamsIsNil)
		} else {
			handler, ok := {{.Action}}Actions[action]
			if !ok {
				httpx.ErrorCtx(ctx, w, types.ErrActionNotFound)
			} else {
				data, err := handler(svcCtx)(ctx, r)
				if err != nil {
					httpx.ErrorCtx(ctx, w, err)
				} else {
					httpx.OkJsonCtx(ctx, w, data)
				}
			}
		}
	}
}