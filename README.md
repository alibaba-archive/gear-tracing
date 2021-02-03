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

  "github.com/opentracing/basictracer-go"
  "github.com/opentracing/opentracing-go"
  "github.com/teambition/gear"
  "github.com/teambition/gear-tracing"
)

func init() {
  opentracing.SetGlobalTracer(basictracer.New(basictracer.NewInMemoryRecorder()))
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
  router := gear.NewRouter()
  router.Use(tracing.New())
  router.Get("/", func(ctx *gear.Context) error {
    span := opentracing.SpanFromContext(ctx)
    assert.NotNil(span)
    span.SetTag("testing", "testing")
    return ctx.HTML(200, "Test Tracing")
  })

  app.UseHandler(router)
  app.Listen(":3000")
}
```
