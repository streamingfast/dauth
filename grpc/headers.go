package grpc

import (
	"context"
	"net/url"

	"github.com/streamingfast/dauth"

	"google.golang.org/grpc/metadata"
)

func convertHeaders(headers url.Values) (md metadata.MD) {
	md = make(map[string][]string)
	for k, v := range headers {
		if len(v) > 0 {
			md[k] = append(md[k], v...)
		}
	}
	return
}

func convertMetadata(md metadata.MD) (headers url.Values) {
	headers = make(map[string][]string)
	for k, v := range md {
		if len(v) > 0 {
			headers[k] = append(headers[k], v...)
		}
	}
	return
}

func HeadersFromContext(ctx context.Context) url.Values {
	md, _ := metadata.FromIncomingContext(ctx)
	return convertMetadata(md)
}

func HeadersToContext(ctx context.Context, headers url.Values) context.Context {
	return metadata.NewOutgoingContext(ctx, convertHeaders(headers))
}

func IdentityFromMetadata(md metadata.MD) dauth.Identity {
	return dauth.Identity{
		UserID:   md.Get(dauth.SFHeaderUserID)[0],
		ApiKeyID: md.Get(dauth.SFHeaderApiKeyID)[0],
		IP:       md.Get(dauth.SFHeaderIP)[0],
	}
}

func IdentityFromContext(ctx context.Context) dauth.Identity {
	md, _ := metadata.FromIncomingContext(ctx)
	return IdentityFromMetadata(md)
}

func IdentityToContext(ctx context.Context, identity dauth.Identity) context.Context {
	md := metadata.Pairs(
		dauth.SFHeaderUserID, identity.UserID,
		dauth.SFHeaderApiKeyID, identity.ApiKeyID,
		dauth.SFHeaderIP, identity.IP,
	)
	return metadata.NewOutgoingContext(ctx, md)
}
