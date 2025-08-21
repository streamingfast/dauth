package dauth

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	HeaderUserID   string = "x-user-id"
	HeaderApiKeyID string = "x-api-key-id"
	HeaderMeta     string = "x-meta"
	HeaderPlanTier string = "x-plan-tier" // As of August 2025, one of FREE, SCALING, PRO, ENTERPRISE
	HeaderIP       string = "x-real-ip"

	deprecatedHeaderUserID   string = "x-sf-user-id"
	deprecatedHeaderApiKeyID string = "x-sf-api-key-id"
	deprecatedHeaderMeta     string = "x-sf-meta"
)

type TrustedHeaders map[string]string

type trustedHeadersKeyType int

const trustedHeadersKey trustedHeadersKeyType = iota

func WithTrustedHeaders(ctx context.Context, h TrustedHeaders) context.Context {
	lowercased := make(TrustedHeaders)
	for k, v := range h {
		lowercased[strings.ToLower(k)] = v
	}

	return context.WithValue(ctx, trustedHeadersKey, lowercased)
}

func FromContext(ctx context.Context) TrustedHeaders {
	val := ctx.Value(trustedHeadersKey)
	if val == nil {
		return nil
	}
	return val.(TrustedHeaders)
}

func (h TrustedHeaders) UserID() string {
	if u, ok := h[HeaderUserID]; ok {
		return u
	}
	return h[deprecatedHeaderUserID]
}

func (h TrustedHeaders) APIKeyID() string {
	if u, ok := h[HeaderApiKeyID]; ok {
		return u
	}
	return h[deprecatedHeaderApiKeyID]
}

func (h TrustedHeaders) Meta() string {
	if u, ok := h[HeaderMeta]; ok {
		return u
	}
	return h[deprecatedHeaderMeta]
}

func (h TrustedHeaders) RealIP() string {
	return h[HeaderIP]
}

func (h TrustedHeaders) PlanTier() string {
	return h[HeaderPlanTier]
}

func (h TrustedHeaders) Get(key string) string {
	return h[strings.ToLower(key)]
}

func (h TrustedHeaders) ToOutgoingGRPCContext(ctx context.Context) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.New(h))
}
