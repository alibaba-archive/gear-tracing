Tracegen
========

Generates interface decorators with [opentracing](http://opentracing.io) support.

Installation
------------

```
go get github.com/teambition/gear-tracing/tracegen
```

Example
-------

```go
type Cache interface {
	Set(ctx context.Context, key, value []byte) error
	Get(ctx context.Context, key []byte) (value []byte, err error)
}
```

```
tracegen -i Cache -o example/cache_trace.go example
```

Will generate:
```go
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

```