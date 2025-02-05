package httpx

import (
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/zeromicro/go-zero/core/syncx"
	"io"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/mapping"
	"github.com/zeromicro/go-zero/rest/internal/encoding"
	"github.com/zeromicro/go-zero/rest/internal/header"
	"github.com/zeromicro/go-zero/rest/pathvar"
)

const (
	formKey           = "form"
	pathKey           = "path"
	maxMemory         = 32 << 20 // 32MB
	maxBodyLen        = 8 << 20  // 8MB
	separator         = ";"
	tokensInAttribute = 2
)

var (
	formUnmarshaler      = mapping.NewUnmarshaler(formKey, mapping.WithStringValues())
	pathUnmarshaler      = mapping.NewUnmarshaler(pathKey, mapping.WithStringValues())
	validatePostJsonBody = syncx.ForAtomicBool(true)
	validate             = validator.New()
)

// SetValidatePostJsonBody 设置是否验证post请求body
func SetValidatePostJsonBody(ok bool) {
	validatePostJsonBody.Set(ok)
}

// Parse parses the request.
func Parse(r *http.Request, v interface{}) error {
	if err := ParsePath(r, v); err != nil {
		return err
	}

	if err := ParseForm(r, v); err != nil {
		return err
	}

	if err := ParseHeaders(r, v); err != nil {
		return err
	}

	return ParseJsonBody(r, v)
}

// ParseHeaders parses the headers request.
func ParseHeaders(r *http.Request, v interface{}) error {
	return encoding.ParseHeaders(r.Header, v)
}

// ParseForm parses the form request.
func ParseForm(r *http.Request, v interface{}) error {
	params, err := GetFormValues(r)
	if err != nil {
		return err
	}

	return formUnmarshaler.Unmarshal(params, v)
}

// ParseHeader parses the request header and returns a map.
func ParseHeader(headerValue string) map[string]string {
	ret := make(map[string]string)
	fields := strings.Split(headerValue, separator)

	for _, field := range fields {
		field = strings.TrimSpace(field)
		if len(field) == 0 {
			continue
		}

		kv := strings.SplitN(field, "=", tokensInAttribute)
		if len(kv) != tokensInAttribute {
			continue
		}

		ret[kv[0]] = kv[1]
	}

	return ret
}

// ParseJsonBody parses the post request which contains json in body.
func ParseJsonBody(r *http.Request, v interface{}) error {
	if withJsonBody(r) {
		reader := io.LimitReader(r.Body, maxBodyLen)

		if err := json.NewDecoder(reader).Decode(v); err != nil {
			return err
		}

		if cv, ok := v.(interface {
			WithContextValue(ctx context.Context)
		}); ok {
			cv.WithContextValue(r.Context())
		}

		if validatePostJsonBody.True() {
			return validate.Struct(v)
		}
		return nil
	}

	return mapping.UnmarshalJsonMap(nil, v)
}

// ParsePath parses the symbols reside in url path.
// Like http://localhost/bag/:name
func ParsePath(r *http.Request, v interface{}) error {
	vars := pathvar.Vars(r)
	m := make(map[string]interface{}, len(vars))
	for k, v := range vars {
		m[k] = v
	}

	return pathUnmarshaler.Unmarshal(m, v)
}

func withJsonBody(r *http.Request) bool {
	return r.ContentLength > 0 && strings.Contains(r.Header.Get(header.ContentType), header.ApplicationJson)
}
