package dauth

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	// Deprecated: Use [HeaderOrganizationID] instead, we might in the future keep both header,
	// but the meaning of both will change. For now, you can assume that user id == organization id
	// as it's how we bill and organize things for now.
	HeaderUserID string = "x-user-id"

	// HeaderOrganizationID is the header carrying the organization id, the actual header is
	// named x-user-id for backward compatibility reasons to keep downstream impact minimal for
	// now, but it's really populated with the organization id.
	HeaderOrganizationID string = "x-user-id"

	HeaderApiKeyID           string = "x-api-key-id"
	HeaderMeta               string = "x-meta"
	HeaderIP                 string = "x-real-ip"
	HeaderSubstreamsPlanTier string = "x-substreams-plan-tier" // As of August 2025, one of FREE, SCALING, PRO, ENTERPRISE

	// Internal: Do not use, only use while transitioning to the new header name,
	// will be removed in the future, always use [HeaderOrganizationID]
	// instead of this.
	HeaderNewOrganizationID string = "x-organization-id"

	deprecatedHeaderUserID     string = "x-user-id"
	deprecatedSfHeaderUserID   string = "x-sf-user-id"
	deprecatedSfHeaderApiKeyID string = "x-sf-api-key-id"
	deprecatedSfHeaderMeta     string = "x-sf-meta"
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

// Deprecated: use [OrganizationID] instead, the [HeaderUserID] now carries the organization id.
// (the [HeaderUserID] is kept for backward compatibility reasons but is also deprecated, use
// [HeaderOrganizationID] instead).
func (h TrustedHeaders) UserID() string {
	if u, ok := h[deprecatedHeaderUserID]; ok {
		return u
	}
	return h[deprecatedSfHeaderUserID]
}

// OrganizationID returns the organization id present in the trusted headers,
// returns "" if not present.
func (h TrustedHeaders) OrganizationID() string {
	if u, ok := h[HeaderNewOrganizationID]; ok {
		return u
	}
	if u, ok := h[deprecatedHeaderUserID]; ok {
		return u
	}
	return h[deprecatedSfHeaderUserID]
}

func (h TrustedHeaders) APIKeyID() string {
	if u, ok := h[HeaderApiKeyID]; ok {
		return u
	}
	return h[deprecatedSfHeaderApiKeyID]
}

func (h TrustedHeaders) Meta() string {
	if u, ok := h[HeaderMeta]; ok {
		return u
	}
	return h[deprecatedSfHeaderMeta]
}

func (h TrustedHeaders) RealIP() string {
	return h[HeaderIP]
}

func (h TrustedHeaders) SubstreamsPlanTier() string {
	return h[HeaderSubstreamsPlanTier]
}

func (h TrustedHeaders) Get(key string) string {
	return h[strings.ToLower(key)]
}

func (h TrustedHeaders) ToOutgoingGRPCContext(ctx context.Context) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.New(h))
}
