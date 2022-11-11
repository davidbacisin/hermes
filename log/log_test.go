package log

import (
	"context"
	"testing"

	"github.com/davidbacisin/hermes/svcmeta"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
)

func Test_Logger_InContext(t *testing.T) {
	t.Run("Embeds logger in the context", func(t *testing.T) {
		ctx := context.Background()
		ctx = svcmeta.Context(ctx, &svcmeta.Resource{
			Name:     "service",
			Version:  "0.0.1",
			Instance: "abcd",
		})
		l1 := logr.Discard().WithName("test")
		ctx = Context(ctx, l1)
		assert.NotNil(t, ctx)

		l2 := FromContext(ctx, "sublogger")
		assert.Equal(t, l1, l2)
	})

	t.Run("Accepts zero-valued parameters", func(t *testing.T) {
		ctx := Context(context.TODO(), logr.Logger{})
		assert.NotNil(t, ctx)

		l := FromContext(ctx, "sublogger")
		assert.NotZero(t, l)
	})

	t.Run("FromContext with uninitialized context", func(t *testing.T) {
		l := FromContext(context.TODO(), "test")
		assert.NotZero(t, l)
	})
}
