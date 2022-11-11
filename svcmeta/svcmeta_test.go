package svcmeta

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Resource_InContext(t *testing.T) {
	t.Run("FromContext without initializing context", func(t *testing.T) {
		ctx := context.Background()
		r := FromContext(ctx)
		assert.NotNil(t, r)
	})

	t.Run("NewInContext", func(t *testing.T) {
		expected := &Resource{
			Name: "abc123",
		}
		ctx := context.Background()
		ctx = Context(ctx, expected)
		actual := FromContext(ctx)

		// Check that it is the same instance of Resource by changing the original
		expected.Version = "2"
		assert.Equal(t, expected, actual)
	})
}

func Test_AsOtelResource(t *testing.T) {
	r := &Resource{
		Name:     "service",
		Version:  "1a49e63f",
		Instance: "1",
	}
	actual := r.AsOtelResource()
	assert.NotNil(t, actual)
}
