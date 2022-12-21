package config

import {{.authImport}}

// 1. 结构体tag default可以设置该字段的默认值
//      使用规则参考：https://github.com/creasty/defaults
// 2. 结构体tag validate可以设置该字段校验
//      使用规则参考：https://github.com/go-playground/validator
//      比如数字范围验证：   Port     int    `validate:"required,gt=0,lte=65535"`
//      比如枚举类型值：     Batcher  string  `default:"jaeger" validate:"oneof=jaeger zipkin grpc"`

type Config struct {
	rest.RestConf
	{{.auth}}
	{{.jwtTrans}}
}
