package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"

	opentracing "github.com/opentracing/opentracing-go"
)

func main() {
	var hostPort string

	if len(os.Args) == 2 {
		hostPort = os.Args[1]
		if !strings.Contains(os.Args[1], ":") {
			hostPort += ":6831"
		}
	}

	// 1) Create a opentracing.Tracer that sends data to Zipkin
	cfg := &config.Configuration{
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  hostPort,
		},
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
	}

	tracer, closer, err := cfg.New(
		"your_service_name",
		config.Logger(jaeger.StdLogger),
	)
	if err != nil {
		panic("Cannot create tracer: " + err.Error())
	}
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
