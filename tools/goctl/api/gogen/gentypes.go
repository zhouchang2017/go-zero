package gogen

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	apiutil "github.com/zeromicro/go-zero/tools/goctl/api/util"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
)

const typesFile = "types"

//go:embed types.tpl
var typesTemplate string

var actionHandler = `type (
	ActionHandler func(ctx context.Context, r *http.Request) (resp interface{}, err error)
)`

var actionErrors = `var (
	// ErrActionParamsIsNil Action参数不存在
	ErrActionParamsIsNil = errors.New("Action params is nil")
	// ErrActionNotFound Action尚未注册
	ErrActionNotFound = errors.New("Action not found")
)`

// BuildTypes gen types to string
func BuildTypes(types []spec.Type, hasActions bool) (string, error) {
	var buf strings.Builder
	if hasActions {
		buf.WriteString(actionHandler)
		buf.WriteByte('\n')
		buf.WriteByte('\n')
	}
	var builder strings.Builder
	first := true
	for _, tp := range types {
		if first {
			first = false
		} else {
			builder.WriteString("\n\n")
		}
		if err := writeType(&builder, tp); err != nil {
			return "", apiutil.WrapErr(err, "Type "+tp.Name()+" generate error")
		}
	}
	buf.WriteString(builder.String())
	return buf.String(), nil
}

func genTypes(dir string, cfg *config.Config, api *spec.ApiSpec) error {
	containsTime := api.ContainsTime

	_, actions, _ := getRoutes(api)
	val, err := BuildTypes(api.Types, len(actions) > 0)
	if err != nil {
		return err
	}

	var builder strings.Builder
	var vars strings.Builder
	if len(actions) > 0 || containsTime {
		builder.WriteByte('\n')
		builder.WriteString("import (\n")
		if len(actions) > 0 {
			vars.WriteString(actionErrors)
			vars.WriteByte('\n')
			builder.WriteString(strconv.Quote("errors"))
			builder.WriteByte('\n')
			builder.WriteString(strconv.Quote("context"))
			builder.WriteByte('\n')
			builder.WriteString(strconv.Quote("net/http"))
			builder.WriteByte('\n')
		}
		if containsTime {
			builder.WriteString(strconv.Quote("time"))
			builder.WriteByte('\n')
		}
		builder.WriteString(")\n")
	}

	typeFilename, err := format.FileNamingFormat(cfg.NamingFormat, typesFile)
	if err != nil {
		return err
	}

	typeFilename = typeFilename + ".go"
	filename := path.Join(dir, typesDir, typeFilename)
	os.Remove(filename)

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          typesDir,
		filename:        typeFilename,
		templateName:    "typesTemplate",
		category:        category,
		templateFile:    typesTemplateFile,
		builtinTemplate: typesTemplate,
		data: map[string]interface{}{
			"importPackages": builder.String(),
			"vars":           vars.String(),
			"types":          val,
		},
	})
}

func writeType(writer io.Writer, tp spec.Type) error {
	structType, ok := tp.(spec.DefineStruct)
	if !ok {
		return fmt.Errorf("unspport struct type: %s", tp.Name())
	}

	fmt.Fprintf(writer, "type %s struct {\n", util.Title(tp.Name()))
	for _, member := range structType.Members {
		if member.IsInline {
			if _, err := fmt.Fprintf(writer, "%s\n", strings.Title(member.Type.Name())); err != nil {
				return err
			}

			continue
		}

		if err := writeProperty(writer, member.Name, member.Tag, member.GetComment(), member.Type, 1); err != nil {
			return err
		}
	}
	fmt.Fprintf(writer, "}")
	return nil
}
