package tracing

import (
	"fmt"

	"github.com/opentracing/opentracing-go"
	"github.com/teambition/gear"
)

// New returns a tracing middleware
func New(opts ...opentracing.StartSpanOption) gear.Middleware {
	return func(ctx *gear.Context) error {
		span := opentracing.StartSpan(fmt.Sprintf(`%s %s`, ctx.Method, ctx.Path), opts...)
		ctx.WithContext(opentracing.ContextWithSpan(ctx.Context(), span))
		ctx.OnEnd(span.Finish)
		return nil
	}
}
