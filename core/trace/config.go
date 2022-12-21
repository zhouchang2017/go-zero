package trace

// TraceName represents the tracing name.
var TraceName = "go-zero"

// A Config is an opentelemetry config.
type Config struct {
	Name     string
	Endpoint string
	Sampler  float64 `default:"1.0"`
	Batcher  string  `default:"jaeger" validate:"oneof=jaeger zipkin grpc"`
}
