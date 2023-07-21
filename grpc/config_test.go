package grpc

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_newConfig(t *testing.T) {
	tests := []struct {
		url       string
		expect    *config
		expectErr bool
	}{
		{
			url: "grpc://localhost:9018",
			expect: &config{
				endpoint:              "localhost:9018",
				enabledContinuousAuth: false,
			},
		},
		{
			url: "grpc://localhost:9018?continuous=true",
			expect: &config{
				endpoint:              "localhost:9018",
				enabledContinuousAuth: true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.url, func(t *testing.T) {
			config, err := newConfig(test.url)
			if test.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expect, config)
			}
		})
	}
}
