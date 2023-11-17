package dauth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestHeadersContext(t *testing.T) {
	ctx := context.Background()

	th := make(TrustedHeaders)

	th[SFHeaderUserID] = "my-user-id"
	th["random-key"] = "102"

	ctx = WithTrustedHeaders(ctx, th)

	got := FromContext(ctx)

	assert.Equal(t, "my-user-id", got.Get(SFHeaderUserID))
	assert.Equal(t, "102", got.Get("random-key"))
}

func TestHeadersMetadataContext(t *testing.T) {
	th := make(TrustedHeaders)
	th[SFHeaderUserID] = "my-user-id"
	th[SFHeaderIP] = "10.2.3.4"

	ctx := context.Background()
	ctx = th.ToOutgoingGRPCContext(ctx)

	md, ok := metadata.FromOutgoingContext(ctx)
	assert.True(t, ok)

	assert.Equal(t, []string{"my-user-id"}, md.Get(SFHeaderUserID))
}

func TestNilDoesntPanic(t *testing.T) {
	ctx := context.Background()

	h := FromContext(ctx)
	assert.Nil(t, h)
	assert.Equal(t, "", h.RealIP())
	assert.Equal(t, "", h.UserID())
	assert.Equal(t, "", h.Meta())
	assert.Equal(t, "", h.Get("something"))
}
