package trust

import (
	"context"
	"strings"

	"github.com/streamingfast/dauth"
)

func Register() {
	dauth.Register("trust", func(configURL string) (dauth.Authenticator, error) {
		return &trustPlugin{}, nil
	})

}

type trustPlugin struct {
}

func (t *trustPlugin) Ready(_ context.Context) bool {
	return true
}

func (t *trustPlugin) Authenticate(ctx context.Context, path string, headers map[string][]string, ipAddress string) (context.Context, error) {
	out := make(dauth.TrustedHeaders)
	for key, values := range headers {
		out[strings.ToLower(key)] = values[0]
	}
	return dauth.WithTrustedHeaders(ctx, out), nil
}
