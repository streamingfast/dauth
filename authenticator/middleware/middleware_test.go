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

package middleware

import (
	"net/http"
	"testing"

	"gotest.tools/assert"
)

func TestExtractTokenFromURL(t *testing.T) {

	tests := []struct {
		name          string
		regex         string
		url           string
		expectedToken string
		expectedPath  string
	}{
		{"vanilla", "^/v2/[^/]*$", "/v2/abcdefg1234567", "abcdefg1234567", "/v2"},
		{"vanilla-fullurl", "^/v2/[^/]*$", "http://example.com/v2/abcdefg1234567", "abcdefg1234567", "/v2"},
		{"matchall", ".*", "/v2/abcdefg1234567", "abcdefg1234567", "/v2"},
		{"longer-fullurl", "^/api/v2/[^/]*$", "http://bob.com/api/v2/abcdefg1234567", "abcdefg1234567", "/api/v2"},
		{"wrongprefix", "^/api/v2/[^/]*$", "/v2/abcdefg1234567", "", "/v2/abcdefg1234567"},
		{"noprefix", "^/api/v2/[^/]*$", "/abcdefg1234567", "", "/abcdefg1234567"},
		{"noprefix", ".*", "/abcdefg1234567", "abcdefg1234567", "/"},
		{"noprefix-notoken", ".*", "/", "", "/"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			opt := WithExtractTokenFromURLPathLastSegment(tc.regex)
			a := &AuthMiddleware{}
			opt(a)
			r, _ := http.NewRequest("", tc.url, nil)
			token := a.tokenExtractFunc(r)
			assert.Equal(t, tc.expectedToken, token)
		})
	}
}
