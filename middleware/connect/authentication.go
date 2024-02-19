package server

import (
	"context"
	"net/http"
	"regexp"

	"github.com/streamingfast/dauth"
	"github.com/streamingfast/dauth/middleware"
	tracing "github.com/streamingfast/sf-tracing"
	"google.golang.org/grpc/metadata"
)

var portSuffixRegex = regexp.MustCompile(`:[0-9]{2,5}$`)
var EmptyMetadata = metadata.New(nil)

func validateAuth(
	ctx context.Context,
	path string,
	headers http.Header,
	peerAddr string,
	authenticator dauth.Authenticator) (context.Context, error) {

	if traceId := tracing.GetTraceID(ctx).String(); traceId != "" {
		headers.Set("x-trace-id", traceId)
	}

	childCtx, err := authenticator.Authenticate(ctx, path, headers, middleware.RealIP(peerAddr, headers))
	if err != nil {
		return nil, err
	}

	return childCtx, nil
}
