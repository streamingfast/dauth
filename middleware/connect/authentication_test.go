package server

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/streamingfast/dauth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		// Lower-case like real authenticators do, otherwise "X-SUBSTREAMS-Ll" and "x-substreams-ll"
		// collide in [dauth.WithTrustedHeaders] and the winner depends on map iteration order.
		out[strings.ToLower(key)] = values[0]
	}
	out["x-substreams-ll"] = "987"
	out["x-user-id"] = "a1b2c3"
	return dauth.WithTrustedHeaders(ctx, out), nil
}

func Test_validAuth(t *testing.T) {
	headers := http.Header{
		"authorization":   []string{"bearer jwt_token"},
		"X-SUBSTREAMS-Ll": []string{"123"},
	}

	authenticator := &testAuthenticators{}

	ctx, err := validateAuth(context.Background(), "/package.service/method", headers, "127.0.0.2", authenticator)
	require.NoError(t, err)
	trusted := dauth.FromContext(ctx)

	assert.Equal(t, "987", trusted.Get("x-substreams-ll"))
	assert.Equal(t, "a1b2c3", trusted.Get("X-User-ID"))

}
