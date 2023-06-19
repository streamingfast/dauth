package server

import (
	"context"
	"github.com/streamingfast/dauth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net/url"
	"regexp"
)

var portSuffixRegex = regexp.MustCompile(`:[0-9]{2,5}$`)
var EmptyMetadata = metadata.New(nil)

type AuthenticatedServerStream struct {
	grpc.ServerStream
	AuthenticatedContext context.Context
}

func (s AuthenticatedServerStream) Context() context.Context {
	return s.AuthenticatedContext
}

func validateAuth(ctx context.Context, path string, authenticator dauth.Authenticator) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = EmptyMetadata
	}

	ctx, authenticatedheaders, err := authenticator.Authenticate(ctx, path, url.Values(md), extractGRPCRealIP(ctx, md))
	if err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, "authentication: %s", err.Error())
	}

	return metadata.NewIncomingContext(ctx, authenticatedheaders), nil
}
