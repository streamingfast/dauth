package http

import (
	"context"
	"google.golang.org/grpc/metadata"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testAuthenticators struct {
}

func (t testAuthenticators) Close() error { return nil }

func (t testAuthenticators) Authenticate(ctx context.Context, path string, headers map[string][]string, ipAddress string) (metadata.MD, error) {
	out := metadata.MD{}
	for key, values := range headers {
		out.Set(key, values...)
	}
	out.Set("X-SF-SUBSTREAMS-LL", "987")
	out.Set("X-Sf-User-Id", "a1b2c3")
	return out, nil
}

func TestAuthMiddleware_validateAuth(t *testing.T) {
	request, err := http.NewRequest("GET", "http://api.example.com/v1/transactions", nil)
	require.NoError(t, err)

	request.Header.Set("authorization", "bearer jwt_token")
	request.Header.Set("X-SF-SUBSTREAMS-LL", "123")

	authenticator := &testAuthenticators{}
	newRequest, err := validateAuth(request, authenticator)
	require.NoError(t, err)

	assert.Equal(t, "a1b2c3", newRequest.Header.Get("X-SF-USER-ID"))
	assert.Equal(t, "987", newRequest.Header.Get("x-sf-substreams-ll"))
}
