package server

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/streamingfast/dauth"
)

func NewAuthInterceptor(check dauth.Authenticator) *AuthInterceptor {
	return &AuthInterceptor{
		check: check,
	}
}

type AuthInterceptor struct {
	check dauth.Authenticator
}

// WrapUnary implements [Interceptor] by applying the interceptor function.
func (i *AuthInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {

		peerAddr := req.Peer().Addr
		headers := req.Header()
		path := req.Spec().Procedure

		childCtx, newHeaders, err := validateAuth(ctx, path, headers, peerAddr, i.check)
		if err != nil {
			return nil, err
		}

		// wipe existing headers, only keep the ones that passthrough auth
		for h := range headers {
			req.Header().Del(h)
		}
		// tweak the headers on the request
		for k, v := range newHeaders {
			first := true
			for _, vv := range v {
				if first {
					req.Header().Set(k, vv)
					first = false
					continue
				}
				req.Header().Add(k, vv)
			}
		}

		return next(childCtx, req)
	}
}

func (i *AuthInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {

		peerAddr := conn.Peer().Addr
		headers := conn.RequestHeader()
		path := conn.Spec().Procedure

		childCtx, newHeaders, err := validateAuth(ctx, path, headers, peerAddr, i.check)
		if err != nil {
			return err
		}

		// wipe existing headers, only keep the ones that passthrough auth
		for h := range headers {
			conn.RequestHeader().Del(h)
		}
		// tweak the headers on the request
		for k, v := range newHeaders {
			first := true
			for _, vv := range v {
				if first {
					conn.RequestHeader().Set(k, vv)
					first = false
					continue
				}
				conn.RequestHeader().Add(k, vv)
			}
		}

		return next(childCtx, conn)
	}
}

// Noop
func (i *AuthInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

// req.Header().Get("Some-Header")
