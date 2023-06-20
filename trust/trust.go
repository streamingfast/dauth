package trust

import (
	"context"

	"github.com/streamingfast/dauth"
	"google.golang.org/grpc/metadata"
)

func Register() {
	dauth.Register("trust", func(configURL string) (dauth.Authenticator, error) {
		return &trustPlugin{}, nil
	})

}

type trustPlugin struct {
}

func (t *trustPlugin) Authenticate(ctx context.Context, path string, headers map[string][]string, ipAddress string) (context.Context, error) {
	out := metadata.MD{}
	for key, values := range headers {
		out.Set(key, values...)
	}
	return metadata.NewIncomingContext(ctx, out), nil
}
