package grpc

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_parseURL(t *testing.T) {

	tests := []struct {
		url       string
		expect    string
		expectErr bool
	}{

		{
			url:       "grpc://localhost9018",
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
