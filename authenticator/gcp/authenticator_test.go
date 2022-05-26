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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseURL(t *testing.T) {
	cfgURL := "cloud-gcp://projects/foo-project/locations/global/keyRings/bar-keyring/cryptoKeys/default/cryptoKeyVersions/1?ip_whitelist=10.5.*.{1,2,5}"

	kmsKeyPath, ipWhitelist, err := parseURL(cfgURL)

	require.NoError(t, err)
	require.Equal(t, "projects/foo-project/locations/global/keyRings/bar-keyring/cryptoKeys/default/cryptoKeyVersions/1", kmsKeyPath)

	assert.True(t, ipWhitelist.Match("10.5.24.5"))
	assert.False(t, ipWhitelist.Match("10.5.24.7"))  // curly braces
	assert.False(t, ipWhitelist.Match("110.5.24.2")) // prefix attack
	assert.False(t, ipWhitelist.Match("10.5.24.54")) // suffix attack
}
