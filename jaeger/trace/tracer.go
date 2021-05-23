package trace

import (
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func Init(svcName string) (opentracing.Tracer, io.Closer, error) {
	cfg := &config.Configuration{
		ServiceName: svcName,

		// "const" sampler is a binary sampling strategy: 0=never sample, 1=always sample.
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},

		// Log the emitted spans to stdout.
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}
	return cfg.NewTracer(config.Logger(jaeger.StdLogger))
}
