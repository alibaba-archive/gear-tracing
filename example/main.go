package main

import (
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/teambition/gear"
	"github.com/teambition/gear-tracing"
)

func init() {
	opentracing.SetGlobalTracer(opentracing.NoopTracer{})
	// use zipkin tracer
	// collector, err := zipkintracer.NewScribeCollector("127.0.0.1:9410", 3*time.Second)
	// if err == nil {
	// 	tracer, err := zipkintracer.NewTracer(zipkintracer.NewRecorder(collector, false, "https://github.com", "gear-tracing"))
	// 	if err == nil {
	// 		opentracing.SetGlobalTracer(tracer)
	// 	}
	// }
}

func main() {
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
