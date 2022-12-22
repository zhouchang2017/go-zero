package trace

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
	"sync"

	"github.com/zeromicro/go-zero/core/lang"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
)

const (
	kindJaeger = "jaeger"
	kindZipkin = "zipkin"
	kindGrpc   = "grpc"
)

var (
	agents = make(map[string]lang.PlaceholderType)
	lock   sync.Mutex
	tp     *sdktrace.TracerProvider
)

//SetTraceName 设置tracer名称；默认去ServiceConf.Name
func SetTraceName(name string) {
	if name != "" {
		TraceName = name
	}
}

// StartAgent starts an opentelemetry agent.
func StartAgent(c Config) {
	lock.Lock()
	defer lock.Unlock()

	_, ok := agents[c.key()]
	if ok {
		return
	}

	// if error happens, let later calls run.
	if err := startAgent(c); err != nil {
		return
	}

	agents[c.key()] = lang.Placeholder
}

// StopAgent shuts down the span processors in the order they were registered.
func StopAgent() {
	_ = tp.Shutdown(context.Background())
}

func createExporter(c Config) (sdktrace.SpanExporter, error) {
	// Just support jaeger and zipkin now, more for later
	switch c.Batcher {
	case kindJaeger:
		if c.Addr != "" {
			split := strings.Split(c.Addr, ":")
			if len(split) == 2 {
				// 通过agent上报，走udp
				return jaeger.New(jaeger.WithAgentEndpoint(
					jaeger.WithAgentHost(split[0]),
					jaeger.WithAgentPort(split[1]),
				))
			} else {
				return nil, fmt.Errorf("addr err: %s", c.Addr)
			}
		}
		return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(c.Endpoint)))
	case kindZipkin:
		return zipkin.New(c.Endpoint)
	case kindGrpc:
		return otlptracegrpc.NewUnstarted(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(c.Endpoint),
			otlptracegrpc.WithDialOption(grpc.WithBlock()),
		), nil
	default:
		return nil, fmt.Errorf("unknown exporter: %s", c.Batcher)
	}
}

func startAgent(c Config) error {
	opts := []sdktrace.TracerProviderOption{
		// Set the sampling rate based on the parent span to 100%
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(c.Sampler))),
		// Record information about this application in a Resource.
		sdktrace.WithResource(resource.NewSchemaless(semconv.ServiceNameKey.String(c.Name))),
	}

	if len(c.Endpoint) > 0 || len(c.Addr) > 0 {
		exp, err := createExporter(c)
		if err != nil {
			logx.GlobalLogger().Error(err)
			return err
		}

		// Always be sure to batch in production.
		opts = append(opts, sdktrace.WithBatcher(exp))
	}

	tp = sdktrace.NewTracerProvider(opts...)
	otel.SetTracerProvider(tp)
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		logx.GlobalLogger().Errorf("[otel] error: %v", err)
	}))

	return nil
}
