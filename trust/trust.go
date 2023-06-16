package trust

import (
	"context"
	"net/url"

	"github.com/streamingfast/dauth"
)

func init() {
	dauth.Register("trust", func(configURL string) (dauth.Authenticator, error) {
		return &trustPlugin{}, nil
	})

	dauth.Register("null", func(configURL string) (dauth.Authenticator, error) {
		return &trustPlugin{}, nil
	})
}

type trustPlugin struct {
}

func (t *trustPlugin) Authenticate(ctx context.Context, path string, headers url.Values, ipAddress string) (url.Values, error) {
	return headers, nil
}
