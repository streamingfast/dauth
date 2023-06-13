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

package http

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_realIPFromHeader(t *testing.T) {

	cases := []struct {
		name          string
		xforwardedFor []string
		remoteAddr    string
		expectedIP    string
	}{
		{
			name:          "sunny path",
			xforwardedFor: []string{"12.34.56.78, 23.45.67.89"},
			expectedIP:    "12.34.56.78",
		},
		{
			name:          "more then 2 ips",
			xforwardedFor: []string{"8.8.8.8,12.34.56.78, 23.45.67.89"},
			expectedIP:    "12.34.56.78",
		},
		{
			name:          "single ip",
			xforwardedFor: []string{"12.34.56.78"},
			expectedIP:    "12.34.56.78",
		},
		{
			name:          "no ip",
			xforwardedFor: []string{""},
			expectedIP:    "0.0.0.0",
		},
		{
			name:          "with junk",
			xforwardedFor: []string{"foo bar, 12.34.56.78, 23.45.67.89"},
			expectedIP:    "12.34.56.78",
		},
		{
			name:       "from remote addr",
			remoteAddr: "12.34.56.78",
			expectedIP: "12.34.56.78",
		},
		{
			name:       "from remote addr with port",
			remoteAddr: "12.34.56.78:54321",
			expectedIP: "12.34.56.78",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := &http.Request{
				Header: map[string][]string{"X-Forwarded-For": c.xforwardedFor},
			}
			req.RemoteAddr = c.remoteAddr
			ip := realIPFromRequest(req)
			assert.Equal(t, c.expectedIP, ip)
		})
	}
}
