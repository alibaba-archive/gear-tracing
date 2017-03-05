# gear-tracing

[![Build Status](http://img.shields.io/travis/teambition/gear-tracing.svg?style=flat-square)](https://travis-ci.org/teambition/gear-tracing)
[![Coverage Status](http://img.shields.io/coveralls/teambition/gear-tracing.svg?style=flat-square)](https://coveralls.io/r/teambition/gear-tracing)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/teambition/gear-tracing/master/LICENSE)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/teambition/gear-tracing)

Opentracing middleware for Gear.

## Demo

```go
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
  //   tracer, err := zipkintracer.NewTracer(zipkintracer.NewRecorder(collector, false, "https://github.com", "gear-tracing"))
  //   if err == nil {
  //     opentracing.SetGlobalTracer(tracer)
  //   }
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
```
