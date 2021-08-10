// Copyright 2019 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/streamingfast/derr"
	"github.com/streamingfast/dauth/authenticator"
	"google.golang.org/grpc/codes"
)

type AuthErrorHandler = func(w http.ResponseWriter, ctx context.Context, err error)
type AuthMiddleware struct {
	errorHandler     AuthErrorHandler
	authenticator    authenticator.Authenticator
	tokenExtractFunc func(*http.Request) string
}

type Option func(*AuthMiddleware)

func NewAuthMiddleware(authenticator authenticator.Authenticator, errorHandler AuthErrorHandler, options ...Option) *AuthMiddleware {
	a := &AuthMiddleware{
		authenticator: authenticator,
		errorHandler:  errorHandler,
	}
	for _, opt := range options {
		opt(a)
	}
	return a
}

func WithCustomTokenExtractor(f func(*http.Request) string) Option {
	return func(a *AuthMiddleware) {
		a.tokenExtractFunc = f
	}
}

func (middleware *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		tokenString := fromQueryOrHeader(r)
		if middleware.tokenExtractFunc != nil {
			if extractedToken := middleware.tokenExtractFunc(r); extractedToken != "" { // always run this function if it exist, to trim request. Extracted token has precedence over token from Query or Header
				tokenString = extractedToken
			}
		}

		tokenRequirement := middleware.authenticator.GetAuthTokenRequirement()
		if tokenRequirement == authenticator.AuthTokenRequired && tokenString == "" {
			middleware.errorHandler(w, ctx, derr.Status(codes.Unauthenticated, "required authorization token not found"))
			return
		}

		ip := authenticator.RealIPFromRequest(r)
		nextCtx, err := middleware.authenticator.Check(ctx, tokenString, ip)
		if err != nil {
			middleware.errorHandler(w, ctx, derr.Statusf(codes.Unauthenticated, "invalid token provided: %s", err.Error()))
			return
		}

		next.ServeHTTP(w, r.WithContext(nextCtx))
	})
}

// fromQueryOrHeader first looks for a token in the HTTP header "Authorization".
//                   if not found, it will look for a query string "token"
func fromQueryOrHeader(r *http.Request) string {
	headerToken, err := fromAuthHeader(r)
	if err != nil || headerToken == "" {
		formToken := r.URL.Query().Get("token")
		if formToken != "" {
			return formToken
		}

		return ""
	}

	return headerToken
}

func fromAuthHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil // No error, just no token
	}

	authHeaderParts := strings.Fields(authHeader)
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", errors.New("Authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}
