package server

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	headers := http.Header{
		"authorization":      []string{"bearer jwt_token"},
		"X-SF-SUBSTREAMS-LL": []string{"123"},
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
