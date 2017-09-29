package tracing

import (
	"fmt"

	"github.com/opentracing/opentracing-go"
	"github.com/teambition/gear"
)

// New returns a tracing middleware
func New(opts ...opentracing.StartSpanOption) gear.Middleware {
	return func(ctx *gear.Context) error {
		// Attempt to join a trace by getting trace context from the headers.
		wireContext, err := opentracing.GlobalTracer().Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(ctx.Req.Header))
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			fmt.Println("failed parsing trace information from header ", err)
		}
		// copy opts avoiding append in the same opts each time.
		// ChildOf will ignore the nil wireContext.
		opts := append([]opentracing.StartSpanOption{opentracing.ChildOf(wireContext)}, opts...)
		span := opentracing.StartSpan(fmt.Sprintf(`%s %s`, ctx.Method, ctx.Path), opts...)
		ctx.WithContext(opentracing.ContextWithSpan(ctx.Context(), span))
		ctx.OnEnd(span.Finish)
		return nil
	}
}
