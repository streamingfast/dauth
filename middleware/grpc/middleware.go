package server

import (
	"context"
	"github.com/streamingfast/dauth"
	"google.golang.org/grpc"
)

func UnaryAuthChecker(check dauth.Authenticator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		childCtx, err := validateAuth(ctx, info.FullMethod, check)
		if err != nil {
			return nil, err
		}

		return handler(childCtx, req)
	}
}

func StreamAuthChecker(check dauth.Authenticator) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		childCtx, err := validateAuth(ss.Context(), info.FullMethod, check)
		if err != nil {
			return err
		}

		return handler(srv, AuthenticatedServerStream{ServerStream: ss, AuthenticatedContext: childCtx})
	}
}
