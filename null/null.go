package trust

import (
	"context"

	"github.com/streamingfast/dauth"
	"go.uber.org/zap"
)

func Register() {
	dauth.Register("null", func(configURL string, _ *zap.Logger) (dauth.Authenticator, error) {
		return &nullPlugin{}, nil
	})
}

type nullPlugin struct {
}

func (t *nullPlugin) Ready(_ context.Context) bool {
	return true
}

func (t *nullPlugin) Authenticate(ctx context.Context, _ string, _ map[string][]string, _ string) (context.Context, error) {
	return ctx, nil
}
