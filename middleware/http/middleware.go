package http

import (
	"context"
	"github.com/streamingfast/dauth"
	"github.com/streamingfast/derr"
	"google.golang.org/grpc/codes"
	"net/http"
	"net/url"
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

		if err := validateAuth(r, m.authenticator); err != nil {
			m.errorHandler(w, ctx, derr.Statusf(codes.Unauthenticated, "authenticate : %s", err.Error()))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func validateAuth(r *http.Request, authenticator dauth.Authenticator) error {
	ctx := r.Context()

	authenticatedHeaders, err := authenticator.Authenticate(ctx, r.URL.String(), url.Values(r.Header), realIPFromRequest(r))
	if err != nil {
		return err
	}

	for key, values := range authenticatedHeaders {
		for _, value := range values {
			r.Header.Set(key, value)
		}
	}
	return nil
}
