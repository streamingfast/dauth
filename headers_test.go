package dauth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
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

func TestAuthHeaders(t *testing.T) {
	tests := []struct {
		name       string
		headers    TrustedHeaders
		method     string
		bucket     string
		iat, exp   int64
		iatOK      bool
		expOK      bool
		wantFields []zap.Field
	}{
		{name: "nil", headers: nil, wantFields: []zap.Field{}},
		{name: "empty", headers: TrustedHeaders{}, wantFields: []zap.Field{}},
		{
			name:       "api_key",
			headers:    TrustedHeaders{HeaderAuthMethod: "api_key"},
			method:     "api_key",
			wantFields: []zap.Field{zap.String("auth_method", "api_key")},
		},
		{
			name: "jwt valid",
			headers: TrustedHeaders{
				HeaderAuthMethod:            "jwt",
				HeaderAuthJWTLifetimeBucket: "le_15m",
				HeaderAuthJWTIssuedAt:       "1700000000",
				HeaderAuthJWTExpiresAt:      "1700000900",
			},
			method: "jwt", bucket: "le_15m",
			iat: 1700000000, iatOK: true, exp: 1700000900, expOK: true,
			wantFields: []zap.Field{
				zap.String("auth_method", "jwt"),
				zap.String("auth_jwt_lifetime_bucket", "le_15m"),
				zap.Int64("auth_jwt_iat", 1700000000),
				zap.Int64("auth_jwt_exp", 1700000900),
			},
		},
		{
			name: "jwt unparseable iat/exp",
			headers: TrustedHeaders{
				HeaderAuthMethod:            "jwt",
				HeaderAuthJWTLifetimeBucket: "unknown",
				HeaderAuthJWTIssuedAt:       "not-a-number",
				HeaderAuthJWTExpiresAt:      "",
			},
			method: "jwt", bucket: "unknown",
			wantFields: []zap.Field{
				zap.String("auth_method", "jwt"),
				zap.String("auth_jwt_lifetime_bucket", "unknown"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.method, tt.headers.AuthMethod())
			assert.Equal(t, tt.bucket, tt.headers.AuthJWTLifetimeBucket())

			iat, ok := tt.headers.AuthJWTIssuedAt()
			assert.Equal(t, tt.iatOK, ok)
			assert.Equal(t, tt.iat, iat)

			exp, ok := tt.headers.AuthJWTExpiresAt()
			assert.Equal(t, tt.expOK, ok)
			assert.Equal(t, tt.exp, exp)

			assert.Equal(t, tt.wantFields, tt.headers.AuthLogFields())
		})
	}
}
