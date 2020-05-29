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

package authenticator

import (
	"context"
)

type authKeyType int

const authKey authKeyType = iota

// WithCredentials creates a child context containing in it the `Credentials` object.
func WithCredentials(ctx context.Context, credentials Credentials) context.Context {
	return context.WithValue(ctx, authKey, credentials)
}

// GetCredentials extracts `Credentials` object from context if it exists, returning it
// if present and `nil` if not found.
func GetCredentials(ctx context.Context) Credentials {
	credentials, ok := ctx.Value(authKey).(Credentials)
	if !ok {
		return &AnonymousCredentials{}
	}

	return credentials
}
