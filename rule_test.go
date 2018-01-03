package tracing

import (
	"net/http"
	"testing"

	"github.com/DavidCai1993/request"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/teambition/gear"
)

func TestRule(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	t.Run("should work", func(t *testing.T) {
		opentracing.SetGlobalTracer(opentracing.NoopTracer{})

		app := gear.New()
		app.Use(NewRule(""))
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
		require.Nil(err)
		assert.Equal(204, res.StatusCode)
	})
}
