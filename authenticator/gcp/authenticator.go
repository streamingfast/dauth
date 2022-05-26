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
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gobwas/glob"
	gcpjwt "github.com/someone1/gcp-jwt-go"
	"github.com/streamingfast/dauth/authenticator"
)

func init() {
	// cloud-gcp://projects/projectname/locations/global/keyRings/keyringname/cryptoKeys/default/cryptoKeyVersions/1
	authenticator.Register("cloud-gcp", func(configURL string) (authenticator.Authenticator, error) {
		kmsKeyPath, ipWhitelist, err := parseURL(configURL)
		if err != nil {
			return nil, fmt.Errorf("cloud-gcp factory: %w", err)
		}

		return newAuthenticator(strings.TrimLeft(kmsKeyPath, "/"), ipWhitelist)
	})
}

func parseURL(configURL string) (kmsKeyPath string, ipWhitelist glob.Glob, err error) {
	urlObject, err := url.Parse(configURL)
	if err != nil {
		return
	}
	kmsKeyPath = urlObject.Host + urlObject.Path

	values := urlObject.Query()
	ipWhitelist, err = glob.Compile(values.Get("ip_whitelist"))

	return
}

type authenticatorPlugin struct {
	kmsVerificationKeyFunc jwt.Keyfunc
	ipWhitelist            glob.Glob
}

func newAuthenticator(kmsKeyPath string, ipWhitelist glob.Glob) (*authenticatorPlugin, error) {
	kmsVerificationKeyFunc, err := gcpjwt.KMSVerfiyKeyfunc(context.Background(), &gcpjwt.KMSConfig{
		KeyPath: kmsKeyPath,
	})
	if err != nil {
		return nil, fmt.Errorf("new kms verify func: %w", err)
	}

	ap := &authenticatorPlugin{
		kmsVerificationKeyFunc: kmsVerificationKeyFunc,
		ipWhitelist:            ipWhitelist,
	}
	return ap, nil
}

func (a *authenticatorPlugin) GetAuthTokenRequirement() authenticator.AuthTokenRequirement {
	return authenticator.AuthTokenRequired
}

func (a *authenticatorPlugin) Check(ctx context.Context, token, ipAddress string) (context.Context, error) {
	var credentials *Credentials

	if ipAddress != "" && a.ipWhitelist.Match(ipAddress) {
		credentials = &Credentials{
			Version: 0,
			IP:      ipAddress,
		}
		return authenticator.WithCredentials(ctx, credentials), nil
	}
	// JWT validation
	credentials = &Credentials{}
	parsedToken, err := jwt.ParseWithClaims(token, credentials, a.kmsVerificationKeyFunc)
	if err != nil {
		return ctx, fmt.Errorf("unable to parse JWT token: %w", err)
	}

	expectedSigningAlgorithm := gcpjwt.SigningMethodKMSES256.Alg()
	actualSigningAlgorithm := parsedToken.Header["alg"]

	if expectedSigningAlgorithm != actualSigningAlgorithm {
		return ctx, fmt.Errorf("invalid JWT token: expected signing method %s, got %s", expectedSigningAlgorithm, actualSigningAlgorithm)
	}

	if !parsedToken.Valid {
		return ctx, errors.New("invalid JWT token: invalid signature")
	}

	credentials.IP = ipAddress
	authContext := authenticator.WithCredentials(ctx, credentials)
	return authContext, nil
}
