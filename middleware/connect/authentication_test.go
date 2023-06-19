package server

import (
	"context"
	"fmt"
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

func Test_validAuth(t *testing.T) {
	headers := http.Header{
		"authorization":      []string{"bearer jwt_token"},
		"X-SF-SUBSTREAMS-Ll": []string{"123"},
	}

	//auth,:= metadata.NewIncomingContext(context.Background(), headers)
	authenticator := &testAuthenticators{}

	ctx := context.Background()
	newHeaders, err := validateAuth(ctx, "/package.service/method", headers, "127.0.0.2", authenticator)
	fmt.Println("new headers are", newHeaders)
	require.NoError(t, err)

	assert.Equal(t, "a1b2c3", newHeaders.Get("X-Sf-User-Id"))
	assert.Equal(t, "987", newHeaders.Get("X-SF-SUBSTREAMS-LL"))
}
