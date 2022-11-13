package trace

import (
	"context"
	"testing"

	"github.com/davidbacisin/hermes/internal/mocks"
	"github.com/davidbacisin/hermes/svcmeta"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func Test_Trace_Context(t *testing.T) {
	testProvider := sdktrace.NewTracerProvider()

	t.Run("Context without metadata", func(t *testing.T) {
		ctx := context.Background()
		ctx = Context(ctx, testProvider)
		assert.NotNil(t, ctx)
	})

	t.Run("FromContext without initializing context", func(t *testing.T) {
		ctx := context.Background()
		tr := FromContext(ctx)
		assert.NotNil(t, tr)
	})

	t.Run("Context with metadata", func(t *testing.T) {
		rsc := &svcmeta.Resource{
			Name: "abc123",
		}
		ctx := context.Background()
		ctx = svcmeta.Context(ctx, rsc)
		ctx = Context(ctx, testProvider)
		tr := FromContext(ctx)
		assert.NotNil(t, tr)
	})
}

func Test_Trace_Start(t *testing.T) {
	t.Run("Start with uninitialized context", func(t *testing.T) {
		ctx := context.Background()
		ctx, s := Start(ctx, "test")
		assert.NotNil(t, ctx)
		assert.NotEmpty(t, s)

		// End and Error must not panic
		s.Error(nil, "")
		s.End()
	})

	t.Run("Start with only trace provider in context", func(t *testing.T) {
		mockSpanProc := mocks.NewSpanProcessor(t)
		mockSpanProc.EXPECT().OnStart(mock.Anything, mock.Anything).Times(1)
		mockSpanProc.EXPECT().OnEnd(mock.Anything).Times(1)

		ctx := context.Background()
		ctx = Context(ctx, sdktrace.NewTracerProvider(
			sdktrace.WithSpanProcessor(mockSpanProc),
		))

		ctx, s := Start(ctx, "test")
		assert.NotNil(t, ctx)
		assert.NotEmpty(t, s)
		s.End()
	})
}

func Test_Trace_Simple(t *testing.T) {
	originalCtx := context.Background()
	actualCtx := originalCtx
	assert.Equal(t, originalCtx, actualCtx)

	end := Simple(&actualCtx, "test")
	assert.NotEqual(t, originalCtx, actualCtx, "Simple should update the context")
	assert.NotNil(t, end)
	end()
}
