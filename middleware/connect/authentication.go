package server

import (
	"context"
	"net/http"
	"regexp"

	"github.com/streamingfast/dauth"
	"github.com/streamingfast/dauth/middleware"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var portSuffixRegex = regexp.MustCompile(`:[0-9]{2,5}$`)
var EmptyMetadata = metadata.New(nil)

func validateAuth(
	ctx context.Context,
	path string,
	headers http.Header,
	peerAddr string,
	authenticator dauth.Authenticator) (context.Context, metadata.MD, error) {

	childCtx, newHeaders, err := authenticator.Authenticate(ctx, path, headers, middleware.RealIP(peerAddr, headers))
	if err != nil {
		return nil, nil, status.Errorf(codes.Unauthenticated, "authentication: %s", err.Error())
	}

	return childCtx, newHeaders, nil
}
