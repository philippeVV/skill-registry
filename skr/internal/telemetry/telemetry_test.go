package telemetry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitNoEndpoint(t *testing.T) {
	// Reset global state
	tracer = nil

	shutdown, err := Init(context.Background(), "", false)
	require.NoError(t, err)
	assert.NotNil(t, shutdown)

	// Tracer should remain nil (no-op mode)
	assert.Nil(t, tracer)

	// Shutdown should be safe to call
	assert.NoError(t, shutdown(context.Background()))
}

func TestInitConsoleExporter(t *testing.T) {
	// Reset global state
	tracer = nil

	shutdown, err := Init(context.Background(), "", true)
	require.NoError(t, err)
	assert.NotNil(t, shutdown)

	// Tracer should be initialized
	assert.NotNil(t, tracer)

	assert.NoError(t, shutdown(context.Background()))

	// Clean up global state
	tracer = nil
}

func TestEmitWithoutInit(t *testing.T) {
	// Reset global state
	tracer = nil

	// All emit functions should be no-ops when tracer is nil
	ctx := context.Background()
	assert.NotPanics(t, func() {
		EmitInstall(ctx, "test-pkg", "1.0.0", "skill", "https://example.com")
	})
	assert.NotPanics(t, func() {
		EmitUninstall(ctx, "test-pkg", "1.0.0")
	})
	assert.NotPanics(t, func() {
		EmitUpdate(ctx, "test-pkg", "1.0.0", "2.0.0")
	})
}
