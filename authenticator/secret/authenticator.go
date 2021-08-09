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
	"net/url"

	"github.com/streamingfast/dauth/authenticator"
	"go.uber.org/zap"
)

func init() {
	// secret://this-is-the-secret-and-fits-in-the-host-field
	authenticator.Register("secret", func(dsn string) (authenticator.Authenticator, error) {
		u, err := url.Parse(dsn)
		if err != nil {
			return nil, err
		}

		secret := u.Host
		if secret == "" {
			return nil, errors.New("missing mandatory env vars config for DFUSE_DAUTH_PLUGIN")
		}

		return newAuthenticator(secret), nil
	})
}

type authenticatorPlugin struct {
	secret string
}

func newAuthenticator(secret string) *authenticatorPlugin {
	if secret == "" {
		panic("Secret cannot be empty string")
	}

	return &authenticatorPlugin{
		secret: secret,
	}
}

func (a *authenticatorPlugin) GetAuthTokenRequirement() authenticator.AuthTokenRequirement {
	return authenticator.AuthTokenRequired
}

func (a *authenticatorPlugin) Check(ctx context.Context, token, ipAddress string) (context.Context, error) {
	if token == a.secret {
		ctx = authenticator.WithCredentials(ctx, newCredentials(ipAddress))
		return ctx, nil
	}

	return ctx, errors.New("invalid authentication token")

}

func (a *authenticatorPlugin) GetLogFields(ctx context.Context) []zap.Field {
	return []zap.Field{
		zap.String("subject", "secret"),
	}
}
