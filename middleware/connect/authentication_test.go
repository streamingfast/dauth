package server

import (
	"context"
	"net/http"
	"testing"

	"github.com/streamingfast/dauth"
	"github.com/stretchr/testify/require"
	"github.com/test-go/testify/assert"
)

type testAuthenticators struct {
}

func (t testAuthenticators) Close() error { return nil }

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

func Test_validAuth(t *testing.T) {
	headers := http.Header{
		"authorization":      []string{"bearer jwt_token"},
		"X-SF-SUBSTREAMS-Ll": []string{"123"},
	}

	authenticator := &testAuthenticators{}

	ctx, err := validateAuth(context.Background(), "/package.service/method", headers, "127.0.0.2", authenticator)
	require.NoError(t, err)
	trusted := dauth.FromContext(ctx)

	assert.Equal(t, "987", trusted.Get("x-sf-substreams-ll"))
	assert.Equal(t, "a1b2c3", trusted.Get("X-Sf-User-ID"))

}
