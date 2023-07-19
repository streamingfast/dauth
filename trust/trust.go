package trust

import (
	"context"
	"net/url"
	"strings"

	"github.com/streamingfast/dauth"
)

func Register() {
	dauth.Register("trust", func(configURL string) (dauth.Authenticator, error) {
		urlObject, err := url.Parse(configURL)
		if err != nil {
			return nil, err
		}

		out := &trustPlugin{}
		if allowed := urlObject.Query().Get("allowed"); allowed != "" {
			out.allowed = make(map[string]bool)
			for _, al := range strings.Split(allowed, ",") {
				out.allowed[strings.ToLower(al)] = true
			}
		}

		return out, nil
	})

}

type trustPlugin struct {
	allowed map[string]bool
}

func (t *trustPlugin) Ready(_ context.Context) bool {
	return true
}

func (t *trustPlugin) Authenticate(ctx context.Context, path string, headers map[string][]string, ipAddress string) (context.Context, error) {
	out := make(dauth.TrustedHeaders)
	for key, values := range headers {
		lowerKey := strings.ToLower(key)
		if t.allowed == nil || t.allowed[lowerKey] {
			out[lowerKey] = values[0]
		}
	}
	return dauth.WithTrustedHeaders(ctx, out), nil
}
