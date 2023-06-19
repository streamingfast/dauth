package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

type testAuthenticators struct {
}

func (t testAuthenticators) Close() error { return nil }

func (t testAuthenticators) Authenticate(ctx context.Context, path string, headers map[string][]string, ipAddress string) (context.Context, metadata.MD, error) {
	out := metadata.MD{}
	for key, values := range headers {
		out.Set(key, values...)
	}
	out.Set("X-SF-SUBSTREAMS-LL", "987")
	out.Set("X-Sf-User-Id", "a1b2c3")
	return ctx, out, nil
}

func Test_validAuth(t *testing.T) {
	headers := metadata.New(map[string]string{
		"authorization":      "bearer jwt_token",
		"X-SF-SUBSTREAMS-LL": "123",
	})

	ctx := metadata.NewIncomingContext(context.Background(), headers)
	authenticator := &testAuthenticators{}

	ctx, err := validateAuth(ctx, "/package.service/method", authenticator)
	require.NoError(t, err)

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = EmptyMetadata
	}
	assert.Equal(t, []string{"a1b2c3"}, md.Get("X-SF-USER-ID"))
	assert.Equal(t, []string{"987"}, md.Get("x-sf-substreams-ll"))
}
