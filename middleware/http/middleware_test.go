package http

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testAuthenticators struct {
}

func (t *testAuthenticators) Close() error { return nil }

func (t *testAuthenticators) Authenticate(ctx context.Context, path string, headers url.Values, ipAddress string) (url.Values, error) {
	headers.Set("X-SF-SUBSTREAMS-LL", "987")
	headers.Set("X-Sf-User-Id", "a1b2c3")
	return headers, nil
}

func TestAuthMiddleware_validateAuth(t *testing.T) {
	request, err := http.NewRequest("GET", "http://api.example.com/v1/transactions", nil)
	require.NoError(t, err)

	request.Header.Set("authorization", "bearer jwt_token")
	request.Header.Set("X-SF-SUBSTREAMS-LL", "123")

	authenticator := &testAuthenticators{}
	require.NoError(t, validateAuth(request, authenticator))

	assert.Equal(t, "a1b2c3", request.Header.Get("X-SF-USER-ID"))
	assert.Equal(t, "987", request.Header.Get("x-sf-substreams-ll"))
}
