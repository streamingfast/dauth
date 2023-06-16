package server

import (
	"context"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

type testAuthenticators struct {
}

func (t *testAuthenticators) Close() error {
	return nil
}

func (t *testAuthenticators) Authenticate(ctx context.Context, path string, headers url.Values, ipAddress string) (url.Values, error) {
	headers.Set("X-SF-SUBSTREAMS-LL", "987")
	headers.Set("X-Sf-User-Id", "a1b2c3")
	return headers, nil
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
