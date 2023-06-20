package dauth

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	SFHeaderUserID   string = "x-sf-user-id"
	SFHeaderApiKeyID string = "x-sf-api-key-id"
	SFHeaderIP       string = "x-real-ip"
)

type TrustedHeaders map[string]string

type trustedHeadersKeyType int

const trustedHeadersKey trustedHeadersKeyType = iota

func WithTrustedHeaders(ctx context.Context, h TrustedHeaders) context.Context {
	return context.WithValue(ctx, trustedHeadersKey, h)
}

func FromContext(ctx context.Context) TrustedHeaders {
	val := ctx.Value(trustedHeadersKey)
	if val == nil {
		return nil
	}
	return val.(TrustedHeaders)
}

func (h TrustedHeaders) UserID() string {
	return h[SFHeaderUserID]
}

func (h TrustedHeaders) APIKeyID() string {
	return h[SFHeaderApiKeyID]
}

func (h TrustedHeaders) RealIP() string {
	return h[SFHeaderIP]
}

func (h TrustedHeaders) Get(key string) string {
	return h[strings.ToLower(key)]
}

func (h TrustedHeaders) ToOutgoingGRPCContext(ctx context.Context) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.New(h))
}
