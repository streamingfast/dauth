package dauth

import (
	"context"
	"testing"

	"github.com/streamingfast/dauth"
	"github.com/streamingfast/dauth/null"
	"github.com/stretchr/testify/require"
)

func TestNullPlugin_NoTrustedHeaderConfigured(t *testing.T) {
	null.Register()

	auth, err := dauth.New("null://", zlog)
	require.NoError(t, err)

	ctx := context.Background()
	require.True(t, auth.Ready(ctx))

	outCtx, err := auth.Authenticate(ctx, "/", map[string][]string{}, "1.1.1.1")
	require.NoError(t, err)

	trustedHeaders := dauth.FromContext(outCtx)
	require.Equal(t, dauth.TrustedHeaders{
		dauth.HeaderIP: "1.1.1.1",
	}, trustedHeaders)
}

func TestNullPlugin_WithTrustedHeaderConfigured(t *testing.T) {
	null.Register()

	auth, err := dauth.New("null://?x-substreams-stage-layer-parallel-executor-max-count=4", zlog)
	require.NoError(t, err)

	ctx := context.Background()
	require.True(t, auth.Ready(ctx))

	outCtx, err := auth.Authenticate(ctx, "/", map[string][]string{}, "1.1.1.1")
	require.NoError(t, err)

	trustedHeaders := dauth.FromContext(outCtx)
	require.Equal(t, dauth.TrustedHeaders{
		dauth.HeaderIP: "1.1.1.1",
		"x-substreams-stage-layer-parallel-executor-max-count": "4",
	}, trustedHeaders)
}
