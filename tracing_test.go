package tracing

import (
	"testing"

	"github.com/DavidCai1993/request"
	"github.com/opentracing/basictracer-go"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.com/teambition/gear"
	"github.com/teambition/gear/middleware/requestid"
)

func TestGearSession(t *testing.T) {
	t.Run("should work", func(t *testing.T) {
		assert := assert.New(t)
		opentracing.SetGlobalTracer(basictracer.New(basictracer.NewInMemoryRecorder()))

		app := gear.New()
		app.Use(requestid.New())

		router := gear.NewRouter()
		router.Use(New())
		router.Get("/", func(ctx *gear.Context) error {
			span := opentracing.SpanFromContext(ctx)
			assert.NotNil(span)
			span.SetTag("testing", "testing")
			return ctx.End(204)
		})

		app.UseHandler(router)
		srv := app.Start()
		defer srv.Close()
		url := "http://" + srv.Addr().String()

		res, err := request.Get(url).End()
		assert.Nil(err)
		assert.Equal(204, res.StatusCode)
	})
}
