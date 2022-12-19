package gogen

import (
	_ "embed"
	"fmt"
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"github.com/zeromicro/go-zero/tools/goctl/vars"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
)

type actionMap struct {
	action string
}

func newActionMap(action string) *actionMap {
	return &actionMap{action: strings.TrimSpace(action)}
}

func (a actionMap) filename() string {
	return a.action + ".go"
}

func (a actionMap) actionTitle() string {
	return strings.Title(a.action)
}

func (a actionMap) registerActionName() string {
	return fmt.Sprintf("Register%sAction", a.actionTitle())
}

func (a actionMap) actionHandlerName() string {
	return fmt.Sprintf("%sHandler", a.actionTitle())
}

func getRoutersByAction(action string, api *spec.ApiSpec) ([]group, error) {
	groups, _, err := getRoutes(api)
	if err != nil {
		return nil, err
	}
	var g []group
	for _, item := range groups {
		if item.actionRouteName == action {
			g = append(g, item)
		}
	}
	return g, nil
}

//go:embed actions.tpl
var actionsTemplate string

func genActions(dir, rootPkg string, cfg *config.Config, api *spec.ApiSpec) error {
	_, actions, err := getRoutes(api)
	if err != nil {
		return err
	}

	for _, action := range actions {
		a := newActionMap(action)
		groups, err := getRoutersByAction(action, api)
		if err != nil {
			return err
		}

		var buf strings.Builder
		var builder strings.Builder
		for _, g := range groups {
			for _, r := range g.routes {
				if r.HandlerDoc != "" {
					builder.WriteString(r.HandlerDoc)
					builder.WriteByte('\n')
				}
				builder.WriteString(fmt.Sprintf("%s(%s, %s)\n",
					a.registerActionName(), strconv.Quote(strings.TrimLeft(r.path, "/")), r.handler))
			}
		}

		if builder.Len() > 0 {
			buf.WriteByte('\n')
			buf.WriteByte('\n')
			buf.WriteString("func init() {\n")
			buf.WriteString(builder.String())
			buf.WriteString("}\n")
		}

		actionFilename := a.filename()
		filename := path.Join(dir, handlerDir, actionFilename)
		os.Remove(filename)

		if err := genFile(fileGenConfig{
			dir:             dir,
			subdir:          handlerDir,
			filename:        actionFilename,
			templateName:    "actionsTemplate",
			category:        category,
			templateFile:    actionsTemplateFile,
			builtinTemplate: actionsTemplate,
			data: map[string]interface{}{
				"PkgName":            "handler",
				"ModuleUrl":          vars.ProjectOpenSourceURL,
				"ImportPackages":     genActionImports(rootPkg, api, action),
				"AddActionInit":      buf.String(),
				"Action":             a.action,
				"ActionTitle":        a.actionTitle(),
				"RegisterActionName": a.registerActionName(),
				"ActionHandlerName":  a.actionHandlerName(),
			},
		}); err != nil {
			return err
		}
	}
	return nil
}

func genActionImports(parentPkg string, api *spec.ApiSpec, actionName string) string {
	importSet := collection.NewSet()
	importSet.AddStr(fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, contextDir)))
	for _, group := range api.Service.Groups {
		action := group.GetAnnotation("action")
		if action != "" && action != actionName {
			continue
		}
		for _, route := range group.Routes {
			folder := route.GetAnnotation(groupProperty)
			if len(folder) == 0 {
				folder = group.GetAnnotation(groupProperty)
				if len(folder) == 0 {
					continue
				}
			}

			importSet.AddStr(fmt.Sprintf("%s \"%s\"", toPrefix(folder),
				pathx.JoinPackages(parentPkg, handlerDir, folder)))
		}
	}
	imports := importSet.KeysStr()
	sort.Strings(imports)
	projectSection := strings.Join(imports, "\n\t")
	return projectSection
}
