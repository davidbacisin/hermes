/*
Package log wraps around github.com/go-logr/logr to adapt it for use with
other Hermes components, such as svcmeta.Resource.
*/
package log

import (
	"context"

	"github.com/davidbacisin/hermes/svcmeta"
	"github.com/go-logr/logr"
)

// Context enhances the provided logger with service metadata
// and adds it to the provided context. It returns the new context with
// the embedded logger, as well as the enhanced logger.
func Context(ctx context.Context, log logr.Logger) context.Context {
	meta := svcmeta.FromContext(ctx)
	if meta != nil {
		if meta.Version != "" {
			log = log.WithValues("version", meta.Version)
		}
		if meta.Instance != "" {
			log = log.WithValues("instance", meta.Instance)
		}
	}
	return logr.NewContext(ctx, log)
}

// FromContext retrieves the logger, if any, from the provided context
// and returns a derived logger with the provided name. If there is
// no logger in the context, a discard logger is returned.
func FromContext(ctx context.Context, name string) logr.Logger {
	l := logr.FromContextOrDiscard(ctx)
	// The zero logger causes panics. Replace it with a discard logger
	if l == (logr.Logger{}) {
		l = logr.Discard()
	}
	return l.WithName(name)
}
