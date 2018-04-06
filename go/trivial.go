package main

import (
	"fmt"
	"os"
	"time"

	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func main() {
	// 1) Create a opentracing.Tracer that sends data to Zipkin
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:	"const",
			Param:	1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:		true,
			BufferFlushInterval:	1 * time.Second,
		},
	}
	tracer, closer, err := cfg.New(
		"your_service_name",
		config.Logger(jaeger.StdLogger),
	)
	defer closer.Close()

	// 2) Demonstrate simple OpenTracing instrumentation
	parent := tracer.StartSpan("Parent")
	for i := 0; i < 20; i++ {
		parent.LogEvent(fmt.Sprintf("Starting child #%d", i))
		child := tracer.StartSpan("Child", opentracing.ChildOf(parent.Context()))
		time.Sleep(50 * time.Millisecond)
		child.Finish()
	}
	parent.LogEvent("A Log")
	parent.Finish()
}
