package server

import (
	"context"
	"net/url"
	"regexp"

	"github.com/streamingfast/dauth"
	"github.com/streamingfast/dauth/middleware"
	tracing "github.com/streamingfast/sf-tracing"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

func validateAuth(ctx context.Context, path string, authenticator dauth.Authenticator, logger *zap.Logger) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = EmptyMetadata
	}

	if traceId := tracing.GetTraceID(ctx).String(); traceId != "" {
		md.Set("x-trace-id", traceId)
	}

	headers := url.Values(md)
	realIP := middleware.RealIP(peerFromContext(ctx), md)

	logger.Debug("validating authentication for gRPC request", zap.String("path", path), zap.String("real_ip", realIP), zap.Any("headers", headers))
	return authenticator.Authenticate(ctx, path, headers, realIP)
}
