package example

/*
This code was automatically generated using github.com/gojuno/generator lib.
			Please DO NOT modify.
*/
import (
	context "context"

	tracing "github.com/teambition/gear-tracing"
)

type CacheTracer struct {
	next   Cache
	prefix string
}

func NewCacheTracer(next Cache, prefix string) *CacheTracer {
	return &CacheTracer{
		next:   next,
		prefix: prefix,
	}
}

func (t *CacheTracer) Get(ctx context.Context, key []byte) (value []byte, err error) {
	span, ctx := tracing.StartSpanFromContext(ctx, t.prefix+".Cache.Get")
	defer span.Finish()

	return t.next.Get(ctx, key)
}

func (t *CacheTracer) Set(ctx context.Context, key []byte, value []byte) error {
	span, ctx := tracing.StartSpanFromContext(ctx, t.prefix+".Cache.Set")
	defer span.Finish()

	return t.next.Set(ctx, key, value)
}
