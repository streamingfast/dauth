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

	ctx = th.ToContext(ctx)

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
