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
	"fmt"
	"net/url"
)

//go:generate go-enum --noprefix -f=$GOFILE --marshal --names

//
// ENUM(
//   AuthTokenRequired
//   AuthTokenOptional
//   AuthTokenDisabled
// )
//
type AuthTokenRequirement uint

type Authenticator interface {
	GetAuthTokenRequirement() AuthTokenRequirement
	Check(ctx context.Context, token, ipAddress string) (context.Context, error)
}

var registry = make(map[string]FactoryFunc)

func New(config string) (Authenticator, error) {
	u, err := url.Parse(config)
	if err != nil {
		return nil, err
	}

	factory := registry[u.Scheme]
	if factory == nil {
		panic(fmt.Sprintf("no Authenticator plugin named \"%s\" is currently registered", u.Scheme))
	}
	return factory(config)
}

type FactoryFunc func(config string) (Authenticator, error)

func Register(name string, factory FactoryFunc) {
	registry[name] = factory
}
