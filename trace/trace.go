package trace

import (
	"context"
	"time"

	"github.com/davidbacisin/hermes/svcmeta"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type tracerContextKeyType struct{}

var (
	Tracer             trace.Tracer
	DurationMetricName string = "duration"
	CounterMetricName  string = "count"

	tracerContextKey tracerContextKeyType
)

type span struct {
	trace.Span
	name      string
	startTime time.Time
}

type Span interface {
	trace.Span
	Error(error, string)
}

func Context(ctx context.Context, provider *sdktrace.TracerProvider) context.Context {
	meta := svcmeta.FromContext(ctx)
	tracer := provider.Tracer(meta.Name)
	return context.WithValue(ctx, tracerContextKey, tracer)
}

func FromContext(ctx context.Context) trace.Tracer {
	if v, ok := ctx.Value(tracerContextKey).(trace.Tracer); ok {
		return v
	}
	meta := svcmeta.FromContext(ctx)
	return otel.Tracer(meta.Name)
}

func Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, Span) {
	ctx, s := FromContext(ctx).Start(ctx, name, opts...)
	return ctx, span{
		Span:      s,
		name:      name,
		startTime: time.Now(),
	}
}

func Simple(ctx *context.Context, name string, opts ...trace.SpanStartOption) func(...trace.SpanEndOption) {
	var s Span
	*ctx, s = Start(*ctx, name, opts...)
	return s.End
}

func (s span) End(opt ...trace.SpanEndOption) {
	s.Span.End(opt...)
}

func (s span) Error(err error, msg string) {
	err = errors.Wrap(err, msg)
	s.SetStatus(codes.Error, msg)
	s.RecordError(err)
}
