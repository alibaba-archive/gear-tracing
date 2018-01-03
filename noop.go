package tracing

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// A NoopTracer is a trivial, minimum overhead implementation of Tracer
// for which all operations are no-ops.
//
// The primary use of this implementation is in libraries, such as RPC
// frameworks, that make tracing an optional feature controlled by the
// end user. A no-op implementation allows said libraries to use it
// as the default Tracer and to write instrumentation that does
// not need to keep checking if the tracer instance is nil.
//
// For the same reason, the NoopTracer is the default "global" tracer
// (see GlobalTracer and SetGlobalTracer functions).
//
// WARNING: NoopTracer does not support baggage propagation.
type NoopTracer struct{}

type noopSpan struct{}
type noopSpanContext struct{}

var (
	defaultNoopSpanContext = noopSpanContext{}
	defaultNoopSpan        = noopSpan{}
	defaultNoopTracer      = NoopTracer{}
)

const (
	emptyString = ""
)

// noopSpanContext:
func (n noopSpanContext) ForeachBaggageItem(handler func(k, v string) bool) {}

// noopSpan:
func (n noopSpan) Context() opentracing.SpanContext                       { return defaultNoopSpanContext }
func (n noopSpan) SetBaggageItem(key, val string) opentracing.Span        { return defaultNoopSpan }
func (n noopSpan) BaggageItem(key string) string                          { return emptyString }
func (n noopSpan) SetTag(key string, value interface{}) opentracing.Span  { return n }
func (n noopSpan) LogFields(fields ...log.Field)                          {}
func (n noopSpan) LogKV(keyVals ...interface{})                           {}
func (n noopSpan) Finish()                                                {}
func (n noopSpan) FinishWithOptions(opts opentracing.FinishOptions)       {}
func (n noopSpan) SetOperationName(operationName string) opentracing.Span { return n }
func (n noopSpan) Tracer() opentracing.Tracer                             { return defaultNoopTracer }
func (n noopSpan) LogEvent(event string)                                  {}
func (n noopSpan) LogEventWithPayload(event string, payload interface{})  {}
func (n noopSpan) Log(data opentracing.LogData)                           {}

// StartSpan belongs to the Tracer interface.
func (n NoopTracer) StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	return defaultNoopSpan
}

// Inject belongs to the Tracer interface.
func (n NoopTracer) Inject(sp opentracing.SpanContext, format interface{}, carrier interface{}) error {
	return nil
}

// Extract belongs to the Tracer interface.
func (n NoopTracer) Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error) {
	return nil, opentracing.ErrSpanContextNotFound
}
