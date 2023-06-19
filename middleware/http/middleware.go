package http

import (
	"context"
	"net/http"

	"github.com/streamingfast/dauth"
	"github.com/streamingfast/dauth/middleware"
	"github.com/streamingfast/derr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
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

	ctx, authenticatedHeaders, err := authenticator.Authenticate(ctx, r.URL.String(), r.Header, middleware.RealIP(r.RemoteAddr, r.Header))
	if err != nil {
		return nil, err
	}

	ctx = metadata.NewIncomingContext(ctx, authenticatedHeaders)
	newRequest := r.Clone(ctx)

	// We cannot simply cast the metadata.MD into an http.Header since they
	// do not format the keys the same (lowercase vs Capitalized)
	httpHeaders := http.Header{}

	for key, values := range authenticatedHeaders {
		for _, v := range values {
			httpHeaders.Set(key, v)
		}

	}

	newRequest.Header = httpHeaders
	return newRequest, nil
}
