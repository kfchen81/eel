package tracing

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"io"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go"
	"time"
	eel_config "github.com/kfchen81/eel/config"
	"github.com/kfchen81/eel/log"
)

var Tracer opentracing.Tracer
var Closer io.Closer

func initJaeger(service string) (opentracing.Tracer, io.Closer) {
	tracingMode := eel_config.ServiceConfig.String("tracing::MODE")
	//
	log.Logger.Infow("[tracing] init open tracing", "tracing_mode", tracingMode)
	var cfg *config.Configuration
	
	if tracingMode == "dev" {
		//cfg = &config.Configuration{
		//	Sampler: &config.SamplerConfig{
		//		Type:  "const",
		//		Param: 1,
		//	},
		//	Reporter: &config.ReporterConfig{
		//		LogSpans: true,
		//		BufferFlushInterval: 1 * time.Second,
		//	},
		//}
		cfg = &config.Configuration{
			Sampler: &config.SamplerConfig{
				Type:  "probabilistic",
				Param: 0.0001,
			},
			Reporter: &config.ReporterConfig{
				LogSpans: false,
				BufferFlushInterval: 5 * time.Second,
			},
		}
	} else {
		cfg = &config.Configuration{
			Sampler: &config.SamplerConfig{
				Type:  "probabilistic",
				Param: 0.4,
			},
			Reporter: &config.ReporterConfig{
				LogSpans: false,
				BufferFlushInterval: 5 * time.Second,
			},
		}
	}
	
	tracer, closer, err := cfg.New(service, config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer
}

func CreateSubSpan(rootSpan opentracing.Span, operationName string) opentracing.Span {
	subSpan := rootSpan.Tracer().StartSpan(
		operationName,
		opentracing.ChildOf(rootSpan.Context()),
	)
	
	return subSpan
}

func init() {
	serviceName := eel_config.ServiceConfig.String("SERVICE_NAME")
	Tracer, Closer = initJaeger(serviceName)
}