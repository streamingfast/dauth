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

package null

import (
	"context"

	"github.com/dfuse-io/dauth/authenticator"
	"go.uber.org/zap"
)

func init() {
	// null://
	authenticator.Register("null", func(dsn string) (authenticator.Authenticator, error) {
		return newAuthenticator(), nil
	})
}

type authenticatorPlugin struct {
}

func newAuthenticator() *authenticatorPlugin {
	return &authenticatorPlugin{}
}

func (a *authenticatorPlugin) IsAuthenticationTokenRequired() bool {
	return false
}

func (a *authenticatorPlugin) Check(ctx context.Context, token, ipAddress string) (context.Context, error) {
	ctx = authenticator.WithCredentials(ctx, newCredentials(ipAddress))
	return ctx, nil

}

func (a *authenticatorPlugin) GetLogFields(ctx context.Context) []zap.Field {
	return []zap.Field{
		zap.String("subject", "null"),
	}
}
