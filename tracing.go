package tracing

import (
	"fmt"

	"github.com/opentracing/opentracing-go"
	"github.com/teambition/gear"
)

// New returns a tracing middleware
func New(opts ...opentracing.StartSpanOption) gear.Middleware {
	return func(ctx *gear.Context) error {
		// copy opts avoiding append in the same opts each time.
		opts := append([]opentracing.StartSpanOption{}, opts...)
		var span opentracing.Span
		opName := fmt.Sprintf(`%s %s`, ctx.Method, ctx.Path)
		// Attempt to join a trace by getting trace context from the headers.
		wireContext, err := opentracing.GlobalTracer().Extract(
			opentracing.TextMap,
			opentracing.HTTPHeadersCarrier(ctx.Req.Header))
		if err != nil {
			// If for whatever reason we can't join, go ahead an start a new root span.
			span = opentracing.StartSpan(opName, opts...)
		} else {
			opts = append(opts, opentracing.ChildOf(wireContext))
			span = opentracing.StartSpan(opName, opentracing.ChildOf(wireContext))
		}

		ctx.WithContext(opentracing.ContextWithSpan(ctx.Context(), span))
		ctx.OnEnd(span.Finish)
		return nil
	}
}
