package named

import (
	"context"
	"time"
)

var globalResolver = newNoopNamedResolver()

func SetGlobalResolver(r Resolver) {
	globalResolver = r
}

func GetGlobalResolver() Resolver {
	if globalResolver == nil {
		SetGlobalResolver(newNoopNamedResolver())
	}
	return globalResolver
}

type (
	// Resolver 名字服务解析器
	Resolver interface {
		GetInstance(ctx context.Context, target string) (Result, error)
	}

	// Result 解析结果
	Result interface {
		// GetEndpoint 实际调用地址
		GetEndpoint() string
		// CallResultReport 调用结果
		CallResultReport(err error, cost time.Duration)
	}
)
