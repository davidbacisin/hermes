/*
Metadata for identifying a service, whether that is a serverless function,
cloud microservice, or client-side application.
*/
package svcmeta

import (
	"context"

	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.6.1"
)

type svcmetaContextKeyType struct{}

var svcmetaContextKey svcmetaContextKeyType

// Resource represents an instance of the running service
type Resource struct {
	// The name of the service
	Name string
	// A version number or commit hash
	Version string
	// Identifier for the currently executing instance of the service, such as
	// node ID, cloud region, IP address, or assigned GUID
	Instance string
}

// Context creates a new context with the Resource embedded as a value.
// The Resource can be retrieved later to link telemetry data with the
// current Resource.
func Context(ctx context.Context, r *Resource) context.Context {
	return context.WithValue(ctx, svcmetaContextKey, r)
}

// FromContext retrieves the Resource from the context, if any, and otherwise
// returns a new, empty resource.
func FromContext(ctx context.Context) *Resource {
	if v, ok := ctx.Value(svcmetaContextKey).(*Resource); ok {
		return v
	}
	return &Resource{}
}

// AsOtelResource creates an OpenTelemetry resource.Resource with the same
// information about the service.
func (r *Resource) AsOtelResource() *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(r.Name),
		semconv.ServiceVersionKey.String(r.Version),
		semconv.ServiceInstanceIDKey.String(r.Instance),
	)
}
