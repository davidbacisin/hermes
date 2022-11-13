package internal

import sdktrace "go.opentelemetry.io/otel/sdk/trace"

type SpanProcessor interface {
	sdktrace.SpanProcessor
}
