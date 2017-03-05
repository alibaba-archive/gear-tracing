package tracing

import (
	"net/http"
	"testing"

	"github.com/DavidCai1993/request"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.com/teambition/gear"
)

func TestGearSession(t *testing.T) {
	t.Run("should work", func(t *testing.T) {
		assert := assert.New(t)
		opentracing.SetGlobalTracer(opentracing.NoopTracer{})

		app := gear.New()
		app.Use(New())
		app.Use(func(ctx *gear.Context) error {
			span, c := opentracing.StartSpanFromContext(ctx, "test_tracing")
			defer span.Finish()

			assert.NotNil(c.Value(http.LocalAddrContextKey))
			return ctx.End(204)
		})

		srv := app.Start()
		defer srv.Close()
		url := "http://" + srv.Addr().String()

		res, err := request.Get(url).End()
		assert.Nil(err)
		assert.Equal(204, res.StatusCode)
	})
}
