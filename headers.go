package dauth

import (
	"context"
	"maps"
	"slices"
	"strings"
	"sync/atomic"

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

type trustedHeadersHolder struct {
	headers atomic.Pointer[TrustedHeaders]
}

// WithTrustedHeaders attaches h to a new context through a fresh holder. Calling it again on a
// context that already carries trusted headers shadows them: the returned context and its
// children no longer see a [ReplaceTrustedHeaders] made on the parent context.
func WithTrustedHeaders(ctx context.Context, h TrustedHeaders) context.Context {
	holder := &trustedHeadersHolder{}
	holder.headers.Store(lowercase(h))

	return context.WithValue(ctx, trustedHeadersKey, holder)
}

// FromContext returns a snapshot of the trusted headers, call it again to see a later [ReplaceTrustedHeaders].
//
// The returned map is shared with every other reader of ctx and must not be modified, a
// concurrent [ReplaceTrustedHeaders] (e.g. from continuous authentication) reads it.
func FromContext(ctx context.Context) TrustedHeaders {
	holder, ok := ctx.Value(trustedHeadersKey).(*trustedHeadersHolder)
	if !ok {
		return nil
	}
	return *holder.headers.Load()
}

// ReplaceTrustedHeaders swaps the trusted headers of ctx and every context derived from it,
// returning false if ctx carries no trusted headers. The swap is atomic but a read-modify-write
// built on [FromContext] is not, callers must ensure a single writer per context.
func ReplaceTrustedHeaders(ctx context.Context, h TrustedHeaders) bool {
	holder, ok := ctx.Value(trustedHeadersKey).(*trustedHeadersHolder)
	if !ok {
		return false
	}
	holder.headers.Store(lowercase(h))
	return true
}

func lowercase(h TrustedHeaders) *TrustedHeaders {
	lowercased := make(TrustedHeaders, len(h))
	for k, v := range h {
		lowercased[strings.ToLower(k)] = v
	}
	return &lowercased
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

// Names returns the list of header names present in the TrustedHeaders.
func (h TrustedHeaders) Names() []string {
	return slices.Collect(maps.Keys(h))
}
