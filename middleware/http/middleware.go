package http

import (
	"context"
	"net/http"

	"github.com/streamingfast/dauth"
	"github.com/streamingfast/dauth/middleware"
	"github.com/streamingfast/derr"
	tracing "github.com/streamingfast/sf-tracing"
	"google.golang.org/grpc/codes"
)

type AuthErrorHandler = func(w http.ResponseWriter, ctx context.Context, err error)
type AuthMiddleware struct {
	errorHandler  AuthErrorHandler
	authenticator dauth.Authenticator
}

type Option func(*AuthMiddleware)

func NewAuthMiddleware(authenticator dauth.Authenticator, errorHandler AuthErrorHandler, options ...Option) *AuthMiddleware {
	a := &AuthMiddleware{
		authenticator: authenticator,
		errorHandler:  errorHandler,
	}
	for _, opt := range options {
		opt(a)
	}
	return a
}

func (m *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}
		ctx := r.Context()
		request, err := validateAuth(r, m.authenticator)
		if err != nil {
			m.errorHandler(w, ctx, derr.Statusf(codes.Unauthenticated, "authenticate : %s", err.Error()))
			return
		}

		next.ServeHTTP(w, request)
	})
}

func validateAuth(r *http.Request, authenticator dauth.Authenticator) (*http.Request, error) {
	ctx := r.Context()

	headers := r.Header.Clone()
	if traceId := tracing.GetTraceID(ctx).String(); traceId != "" {
		headers.Set("x-trace-id", traceId)
	}

	ctx, err := authenticator.Authenticate(ctx, r.URL.String(), headers, middleware.RealIP(r.RemoteAddr, headers))
	if err != nil {
		return nil, err
	}

	return r.Clone(ctx), nil
}
