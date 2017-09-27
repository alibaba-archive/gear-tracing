package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go-opentracing"
)

var (
	zipkin = flag.String("zipkin", "127.0.0.1:9410", "zipkin service")
	method = flag.String("method", "GET", "http method")
	url    = flag.String("url", "http://example.com", "http url")
)

func main() {
	flag.Parse()
	collector, err := zipkintracer.NewScribeCollector(*zipkin, 3*time.Second)
	if err != nil {
		panic(err)
	}
	tracer, err := zipkintracer.NewTracer(
		zipkintracer.NewRecorder(collector, true, "test", "test"),
	)
	if err != nil {
		panic(err)
	}
	opentracing.SetGlobalTracer(tracer)
	req, err := http.NewRequest(*method, *url, nil)
	if err != nil {
		panic(err)
	}
	span := opentracing.StartSpan(
		*method + ":" + *url,
	)
	defer span.Finish()
	opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)
	fmt.Println(http.DefaultClient.Do(req))
}
