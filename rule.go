package tracing

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/teambition/gear"
)

type kv map[string]*regexp.Regexp

func (k kv) set(key, value string) error {
	reg, err := regexp.Compile(value)
	if err != nil {
		return err
	}
	k[key] = reg
	return nil
}

func (k kv) remove(key string) {
	delete(k, key)
}

func (k kv) match(key, value string) bool {
	reg, ok := k[key]
	if !ok {
		return false
	}
	return reg.MatchString(value)
}

func (k kv) find(value string) bool {
	for _, reg := range k {
		if reg.MatchString(value) {
			return true
		}
	}
	return false
}

const (
	header = "header"
	query  = "query"
	path   = "path"
	all    = "all"
)

// rule tracing with rules
type rule struct {
	sync.RWMutex
	headers, querys, paths kv
}

func (r *rule) set(t, key, value string) error {
	if key == "" && t != path {
		return errors.New("must not empty key")
	}
	if value == "" {
		return errors.New("must not empty value")
	}
	r.Lock()
	defer r.Unlock()
	switch t {
	case header:
		return r.headers.set(key, value)
	case query:
		return r.querys.set(key, value)
	case path:
		return r.paths.set(key, value)
	default:
		return fmt.Errorf("unexpected type(%s) only in (header/query/path)", t)
	}
}

func (r *rule) remove(t, key string) {
	r.Lock()
	defer r.Unlock()
	switch t {
	case header:
		r.headers.remove(key)
	case query:
		r.querys.remove(key)
	case path:
		r.paths.remove(key)
	}
}

func (r *rule) list(t string) interface{} {
	r.RLock()
	defer r.RUnlock()

	values := kv{}
	switch t {
	case header:
		values = r.headers
	case query:
		values = r.querys
	case path:
		values = r.paths
	default:
	}

	if t == path {
		result := []string{}
		for _, v := range values {
			result = append(result, v.String())
		}
		return result
	}

	result := map[string]string{}
	for k, v := range values {
		result[k] = v.String()
	}
	return result
}

func (r *rule) matchHeader(headers http.Header) bool {
	r.RLock()
	defer r.RUnlock()

	for key := range headers {
		if r.headers.match(key, headers.Get(key)) {
			return true
		}
	}
	return false
}

func (r *rule) matchQuery(querys url.Values) bool {
	r.RLock()
	defer r.RUnlock()

	for key := range querys {
		if r.querys.match(key, querys.Get(key)) {
			return true
		}
	}
	return false
}

func (r *rule) matchPath(path string) bool {
	r.RLock()
	defer r.RUnlock()

	return r.paths.find(path)
}

var globalRule = &rule{
	headers: kv{},
	querys:  kv{},
	paths:   kv{},
}

// NewRule middleware with rules
func NewRule(prefix string, opts ...opentracing.StartSpanOption) gear.Middleware {
	// tracing url /prefix/tracing/:_type(header,query,path)
	tracingPrefix := fmt.Sprintf("%s/tracing", prefix)
	//  tracing middleware
	mw := New(opts...)
	return func(ctx *gear.Context) error {
		// tracing operations
		if strings.HasPrefix(ctx.Path, tracingPrefix) {
			ps := strings.Split(ctx.Path, "/")
			_type := ps[len(ps)-1]
			switch ctx.Method {
			case http.MethodPost:
				key := ctx.Query("key")
				value := ctx.Query("value")
				if err := globalRule.set(_type, key, value); err != nil {
					ctx.HTML(http.StatusBadRequest, err.Error())
				} else {
					ctx.HTML(http.StatusOK, "success")
				}
			case http.MethodDelete:
				key := ctx.Query("key")
				globalRule.remove(_type, key)
				ctx.HTML(http.StatusOK, "success")
			case http.MethodGet:
				result := globalRule.list(_type)
				ctx.JSON(http.StatusOK, result)
			default:
				ctx.HTML(http.StatusBadRequest, "only support methods (post/delete/get)")
			}
		}
		// match and tracing
		if globalRule.matchPath(ctx.Path) ||
			globalRule.matchHeader(ctx.Req.Header) ||
			globalRule.matchQuery(ctx.Req.URL.Query()) {
			mw(ctx)
		}
		return nil
	}
}

// StartSpanFromContext starts and returns a Span with `operationName`, using
// any Span found within `ctx` as a ChildOfRef. If no such parent could be
// found, StartSpanFromContext creates a noop Span.
//
// The second return value is a context.Context object built around the
// returned Span.
//
// Example usage:
//
//    SomeFunction(ctx context.Context, ...) {
//        sp, ctx := opentracing.StartSpanFromContext(ctx, "SomeFunction")
//        defer sp.Finish()
//        ...
//    }
func StartSpanFromContext(ctx context.Context, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	return startSpanFromContextWithTracer(ctx, opentracing.GlobalTracer(), operationName, opts...)
}

// startSpanFromContextWithTracer is factored out for testing purposes.
func startSpanFromContextWithTracer(ctx context.Context, tracer opentracing.Tracer, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	var span opentracing.Span
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
		span = tracer.StartSpan(operationName, opts...)
	} else {
		span = defaultNoopTracer.StartSpan(operationName, opts...)
	}
	return span, opentracing.ContextWithSpan(ctx, span)
}
