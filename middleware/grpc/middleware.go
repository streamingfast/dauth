package server

import (
	"context"

	"github.com/streamingfast/dauth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryAuthChecker(check dauth.Authenticator, logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		childCtx, err := validateAuth(ctx, info.FullMethod, check, logger)
		if err != nil {
			return nil, obfuscateErrorMessage(err, logger)
		}

		return handler(childCtx, req)
	}
}

func StreamAuthChecker(check dauth.Authenticator, logger *zap.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		childCtx, err := validateAuth(ss.Context(), info.FullMethod, check, logger)
		if err != nil {
			return obfuscateErrorMessage(err, logger)
		}

		return handler(srv, AuthenticatedServerStream{ServerStream: ss, AuthenticatedContext: childCtx})
	}
}

func obfuscateErrorMessage(err error, logger *zap.Logger) error {
	if st, ok := status.FromError(err); ok {
		msg := st.Message()
		switch st.Code() {
		case codes.Internal, codes.Unavailable, codes.Unknown:
			logger.Error("authentication service via gRPC middleware fatal error", zap.Error(err))
			msg = "error with authentication service, please try again later"
		}
		return status.Error(st.Code(), msg)
	} else {
		// A non-gRPC error here means the authenticator rejected the request itself (missing/malformed
		// credentials, expired token, unknown API key). That's a client-side condition, not a server
		// fault, so it must not pollute error logs nor error-rate alerting.
		logger.Debug("authentication service via gRPC middleware non-gRPC error", zap.Error(err))
	}

	return status.Errorf(codes.Unauthenticated, "authentication: %s", err.Error())
}
