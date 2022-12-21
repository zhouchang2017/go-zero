package named

import (
	"context"
	"time"
)

var _ Resolver = noopNamedResolver{}

type (
	noopNamedResolver struct {
	}
	noopNamedResolverResult struct{ addr string }
)

func newNoopNamedResolver() Resolver {
	return noopNamedResolver{}
}

func (n noopNamedResolverResult) GetEndpoint() string {
	return n.addr
}

func (n noopNamedResolverResult) CallResultReport(err error, cost time.Duration) {
	// nothing to do
}

func (n noopNamedResolver) GetInstance(ctx context.Context, target string) (Result, error) {
	return noopNamedResolverResult{addr: target}, nil
}
