package server

import (
	"context"
	"net/http"
	"regexp"

	"github.com/streamingfast/dauth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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

func validateAuth(
	ctx context.Context,
	path string,
	headers http.Header,
	peerAddr string,
	authenticator dauth.Authenticator) (*http.Request, error) {

	_, authenticatedheaders, err := authenticator.Authenticate(ctx, path, headers, extractGRPCRealIP(peerAddr, headers))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication: %s", err.Error())
	}
	ctx = metadata.NewIncomingContext(ctx, authenticatedHeaders)

	return authenticatedheaders, nil
}
