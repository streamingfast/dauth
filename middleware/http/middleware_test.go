package http

import (
	"context"
	"net/http"
	"testing"

	"github.com/streamingfast/dauth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testAuthenticators struct {
}

func (t testAuthenticators) Ready(_ context.Context) bool {
	return true
}

func (t testAuthenticators) Authenticate(ctx context.Context, path string, headers map[string][]string, ipAddress string) (context.Context, error) {
	out := make(dauth.TrustedHeaders)
	for key, values := range headers {
		out[key] = values[0]
	}
	out["x-sf-substreams-ll"] = "987"
	out["x-sf-user-id"] = "a1b2c3"

	return dauth.WithTrustedHeaders(ctx, out), nil
}

func TestAuthMiddleware_validateAuth(t *testing.T) {
	request, err := http.NewRequest("GET", "http://api.example.com/v1/transactions", nil)
	require.NoError(t, err)

	request.Header.Set("authorization", "bearer jwt_token")
	request.Header.Set("X-SF-SUBSTREAMS-LL", "123")

	authenticator := &testAuthenticators{}
	newRequest, err := validateAuth(request, authenticator)
	require.NoError(t, err)

	trusted := dauth.FromContext(newRequest.Context())

	assert.Equal(t, "987", trusted.Get("x-sf-substreams-ll"))
	assert.Equal(t, "a1b2c3", trusted.Get("X-Sf-User-ID"))
}
