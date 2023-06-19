package server

import (
	"context"
	"net/http"
	"testing"

	"google.golang.org/grpc/metadata"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	return metadata.NewIncomingContext(ctx, out), out, nil
}

func Test_validAuth(t *testing.T) {
	headers := http.Header{
		"authorization":      []string{"bearer jwt_token"},
		"X-SF-SUBSTREAMS-Ll": []string{"123"},
	}

	//auth,:= metadata.NewIncomingContext(context.Background(), headers)
	authenticator := &testAuthenticators{}

	ctx, newHeaders, err := validateAuth(context.Background(), "/package.service/method", headers, "127.0.0.2", authenticator)
	require.NoError(t, err)

	assert.Equal(t, []string{"a1b2c3"}, newHeaders.Get("X-SF-USER-ID"))
	assert.Equal(t, []string{"987"}, newHeaders.Get("x-sf-substreams-ll"))

	md, found := metadata.FromIncomingContext(ctx)
	assert.True(t, found)

	assert.Equal(t, []string{"a1b2c3"}, md.Get("X-SF-USER-ID"))
	assert.Equal(t, []string{"987"}, md.Get("x-sf-substreams-ll"))
}
