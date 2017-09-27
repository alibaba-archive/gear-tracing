package main

import (
	"flag"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go-opentracing"
	"github.com/teambition/gear"
	"github.com/teambition/gear-tracing"
)

var (
	zipkin = flag.String("zipkin", "127.0.0.1:9410", "zipkin scribe service address")
)

func main() {
	flag.Parse()
	collector, err := zipkintracer.NewScribeCollector(*zipkin, 3*time.Second)
	if err == nil {
		tracer, err := zipkintracer.NewTracer(zipkintracer.NewRecorder(collector, false, "https://github.com", "gear-tracing"))
		if err == nil {
			opentracing.SetGlobalTracer(tracer)
		}
	}

	app := gear.New()

	app.Use(tracing.New())
	app.Use(func(ctx *gear.Context) error {
		span, _ := opentracing.StartSpanFromContext(ctx, "test_tracing")
		defer span.Finish()

		time.Sleep(time.Second)
		return ctx.HTML(200, "Test Tracing")
	})
	app.Listen(":3000")
}
