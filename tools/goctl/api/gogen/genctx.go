package gogen

import (
	_ "embed"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
)

const ctxFilename = "ctx"

//go:embed ctx.tpl
var ctxTemplate string

func genCtx(dir, rootPkg string, cfg *config.Config, api *spec.ApiSpec) error {
	filename, err := format.FileNamingFormat(cfg.NamingFormat, ctxFilename)
	if err != nil {
		return err
	}

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          contextDir,
		filename:        filename + ".go",
		templateName:    "ctxTemplate",
		category:        category,
		templateFile:    ctxTemplateFile,
		builtinTemplate: ctxTemplate,
		data:            map[string]string{},
	})
}
