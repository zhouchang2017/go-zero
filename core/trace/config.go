package trace

// TraceName represents the tracing name.
var TraceName = "go-zero"

// A Config is an opentelemetry config.
type Config struct {
	Name     string
	Endpoint string
	Addr     string  // 这里不为空字符的话，就走jaeger的Agent，而不是Collector
	Sampler  float64 `default:"1.0"`
	Batcher  string  `default:"jaeger" validate:"oneof=jaeger zipkin grpc"`
}

func (c Config) key() string {
	if c.Addr != "" {
		return c.Addr
	}
	return c.Endpoint
}
