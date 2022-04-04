// Copyright 2021 dfuse Platform Inc.
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

package gcp

import (
	"strings"

	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
)

type Credentials struct {
	jwt.StandardClaims

	// From JWT
	Version      int     `json:"v"`
	Usage        string  `json:"usg"`
	APIKeyID     string  `json:"aki"`
	Origin       string  `json:"origin,omitempty"`
	FeatureFlags []int32 `json:"opts,omitempty"`

	IP string `json:"-"`
}

func (c *Credentials) GetUserID() string {
	userID := c.Subject
	return strings.TrimPrefix(userID, "uid:")
}

func (c *Credentials) GetLogFields() []zap.Field {
	return []zap.Field{
		zap.String("subject", c.Subject),
		zap.String("jti", c.Id),
		zap.String("api_key_id", c.APIKeyID),
		zap.String("ip", c.IP),
	}
}
