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
	"github.com/streamingfast/dauth/authenticator"
	"go.uber.org/zap"
)

var _ authenticator.Credentials = (*credentials)(nil)

type credentials struct {
	ipAddress string
}

func (c *credentials) Features() *authenticator.Features {
	return &authenticator.Features{}
}

func (c *credentials) Identification() *authenticator.Identification {
	return &authenticator.Identification{}
}

func newCredentials(ipAddress string) *credentials {
	return &credentials{
		ipAddress: ipAddress,
	}
}

func (c *credentials) GetLogFields() []zap.Field {
	return []zap.Field{
		zap.String("ip", c.ipAddress),
	}
}
