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

package secret

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/streamingfast/dauth"
	"go.uber.org/zap"
)

func Register() {
	// secret://this-is-the-secret-and-fits-in-the-host-field?[user_id=<value>]&[api_key_id=<value>]
	dauth.Register("secret", func(configURL string, logger *zap.Logger) (dauth.Authenticator, error) {
		authenticator, err := newAuthenticatorFromURL(configURL)
		if err != nil {
			return nil, fmt.Errorf("failed to setup secret config: %w", err)
		}

		logger.Info("setting up secret authenticator",
			zap.String("secret", maskSecret(authenticator.secret)),
		)

		return authenticator, nil
	})
}

var _ dauth.Authenticator = (*authenticator)(nil)

type authenticator struct {
	secret   string
	userID   string
	apiKeyID string
}

func newAuthenticatorFromURL(urlRaw string) (*authenticator, error) {
	urlObject, err := url.Parse(urlRaw)
	if err != nil {
		return nil, fmt.Errorf("unable to parse url: %w", err)
	}

	params := urlObject.Query()

	return &authenticator{
		secret:   urlObject.Host,
		userID:   params.Get("user_id"),
		apiKeyID: params.Get("api_key_id"),
	}, nil
}

func newAuthenticator(secret string, userID string, apiKeyID string) (*authenticator, error) {
	if secret == "" {
		panic("Secret cannot be empty string")
	}

	if secret == "" {
		return nil, errors.New("missing mandatory secret value")
	}

	return &authenticator{
		secret: secret,
	}, nil
}

// Authenticate implements dauth.Authenticator.
func (a *authenticator) Authenticate(ctx context.Context, path string, headers map[string][]string, ipAddress string) (context.Context, error) {
	authorizationHeaders := headers["authorization"]
	if len(authorizationHeaders) == 0 {
		return ctx, errors.New("missing authorization header")
	}

	var lastAuthErr error
	for _, authorizationHeader := range authorizationHeaders {
		if err := a.validateAuthHeader(authorizationHeader); err != nil {
			lastAuthErr = err
		}
	}

	if lastAuthErr != nil {
		return ctx, lastAuthErr
	}

	out := make(dauth.TrustedHeaders)
	out[dauth.SFHeaderIP] = ipAddress
	out[dauth.SFHeaderUserID] = a.userID
	out[dauth.SFHeaderApiKeyID] = a.apiKeyID

	return dauth.WithTrustedHeaders(ctx, out), nil
}

// Ready implements dauth.Authenticator.
func (*authenticator) Ready(context.Context) bool {
	return true
}

func (a *authenticator) validateAuthHeader(value string) error {
	authHeaderParts := strings.Fields(value)

	var token string
	switch len(authHeaderParts) {
	case 1:
		token = authHeaderParts[0]
	case 2:
		if strings.ToLower(authHeaderParts[0]) != "bearer" {
			return fmt.Errorf("authorization header format must be Bearer {token}")
		}

		token = authHeaderParts[1]
	default:
		return fmt.Errorf("authorization header format must be Bearer {token}")
	}

	if token != a.secret {
		return errors.New("invalid authentication token received")
	}

	return nil
}

func maskSecret(in string) string {
	if len(in) < 9 {
		return strings.Repeat("*", len(in))
	}

	return in[:3] + strings.Repeat("*", len(in)-6) + in[len(in)-3:]
}
