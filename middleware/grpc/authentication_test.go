package server

import (
	"context"
	"testing"

	"github.com/streamingfast/dauth"
	"github.com/streamingfast/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

var zlogTest, _ = logging.PackageLogger("test", "github.com/streamingfast/dauth/middleware/grpc/authentication/test")

func init() {
	logging.InstantiateLoggers()
}

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

	out["x-substreams-ll"] = "987"
	out["x-user-id"] = "a1b2c3"
	return dauth.WithTrustedHeaders(ctx, out), nil
}

func Test_validAuth(t *testing.T) {
	headers := metadata.New(map[string]string{
		"authorization":   "bearer jwt_token",
		"X-SUBSTREAMS-LL": "123",
	})

	ctx := metadata.NewIncomingContext(context.Background(), headers)
	authenticator := &testAuthenticators{}

	ctx, err := validateAuth(ctx, "/package.service/method", authenticator, zlogTest)
	require.NoError(t, err)

	trusted := dauth.FromContext(ctx)

	assert.Equal(t, "987", trusted.Get("x-substreams-ll"))
	assert.Equal(t, "a1b2c3", trusted.Get("X-User-ID"))
}
