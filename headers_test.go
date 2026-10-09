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

	th[HeaderUserID] = "my-user-id"
	th["random-key"] = "102"

	ctx = WithTrustedHeaders(ctx, th)

	got := FromContext(ctx)

	assert.Equal(t, "my-user-id", got.Get(HeaderUserID))
	assert.Equal(t, "102", got.Get("random-key"))
}

func TestHeadersMetadataContext(t *testing.T) {
	th := make(TrustedHeaders)
	th[HeaderUserID] = "my-user-id"
	th[HeaderIP] = "10.2.3.4"

	ctx := context.Background()
	ctx = th.ToOutgoingGRPCContext(ctx)

	md, ok := metadata.FromOutgoingContext(ctx)
	assert.True(t, ok)

	assert.Equal(t, []string{"my-user-id"}, md.Get(HeaderUserID))
}

func TestNilDoesntPanic(t *testing.T) {
	ctx := context.Background()

	h := FromContext(ctx)
	assert.Nil(t, h)
	assert.Equal(t, "", h.RealIP())
	assert.Equal(t, "", h.UserID())
	assert.Equal(t, "", h.Meta())
	assert.Equal(t, "", h.SubstreamsPlanTier())
	assert.Equal(t, "", h.Get("something"))
}

func TestReplaceTrustedHeaders(t *testing.T) {
	ctx := WithTrustedHeaders(context.Background(), TrustedHeaders{HeaderUserID: "my-user-id"})
	derived, cancel := context.WithCancel(ctx)
	defer cancel()

	before := FromContext(derived)

	assert.True(t, ReplaceTrustedHeaders(ctx, TrustedHeaders{HeaderUserID: "my-user-id", "X-Meta": "some-meta"}))

	assert.Equal(t, "some-meta", FromContext(derived).Meta(), "derived contexts see the replaced headers")
	assert.Equal(t, "", before.Meta(), "a snapshot taken before the replace is unchanged")
}

func TestReplaceTrustedHeaders_NoHeaders(t *testing.T) {
	assert.False(t, ReplaceTrustedHeaders(context.Background(), TrustedHeaders{HeaderUserID: "my-user-id"}))
	assert.Nil(t, FromContext(context.Background()))
}
