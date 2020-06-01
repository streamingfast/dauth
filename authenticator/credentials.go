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

import "go.uber.org/zap"

type Credentials interface {
	GetLogFields() []zap.Field
	GetUserID() string
}

type AnonymousCredentials struct {
	// MUST implement Credentials
	userID string
	ip string
}

func newAnonymousCredentials() *AnonymousCredentials {
	return &AnonymousCredentials{
		userID: "anonymous",
		ip: "0.0.0.0",
	}
}

func (c *AnonymousCredentials) GetUserID() string {
	return c.userID
}

func (c *AnonymousCredentials) GetLogFields() []zap.Field {
	return []zap.Field{
		zap.String("subject", c.userID),
		zap.String("api_key_id", c.userID),
		zap.String("ip", c.ip),
	}
}
