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
	"flag"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go-opentracing"
	"github.com/teambition/gear"
	"github.com/teambition/gear-tracing"
)

var (
	zipkin = flag.String("zipkin", "127.0.0.1:9410", "zipkin scribe service address")
	rule   = flag.Bool("rule", false, "use tracing rule")
)

func main() {
	flag.Parse()
	collector, err := zipkintracer.NewScribeCollector(*zipkin, 3*time.Second)
	if err != nil {
		panic(nil)
	}
	tracer, err := zipkintracer.NewTracer(
		zipkintracer.NewRecorder(collector, false, "https://github.com", "gear-tracing"),
	)
	if err != nil {
		panic(err)
	}
	opentracing.SetGlobalTracer(tracer)

	app := gear.New()

	if *rule {
		app.Use(tracing.NewRule(""))
	} else {
		app.Use(tracing.New())
	}
	app.Use(func(ctx *gear.Context) error {
		span, _ := tracing.StartSpanFromContext(ctx, "test_tracing")
		defer span.Finish()

		time.Sleep(time.Second)
		return ctx.HTML(200, "Test Tracing")
	})
	app.Listen(":3000")
}
```

[Use Tracegen for code generator](tracegen/README.md)

## Rule apis

```
# path rules
GET http://example.com/{prefix}/tracing/path
POST http://example.com/{prefix}/tracing/path
DELETE http://example.com/{prefix}/tracing/path

# query rules
GET http://example.com/{prefix}/tracing/query
POST http://example.com/{prefix}/tracing/query
DELETE http://example.com/{prefix}/tracing/query

# header rules
GET http://example.com/{prefix}/tracing/header
POST http://example.com/{prefix}/tracing/header
DELETE http://example.com/{prefix}/tracing/header
```