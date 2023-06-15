package grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseURL(t *testing.T) {
	tests := []struct {
		url       string
		expect    string
		expectErr bool
	}{
		{
			url:       "grpc://localhost:9018",
			expect:    "localhost:9018",
			expectErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.url, func(t *testing.T) {

			v, err := parseURL(test.url)
			if test.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expect, v)
			}

		})
	}

}
