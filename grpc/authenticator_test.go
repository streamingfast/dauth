package grpc

import (
	"context"
	"fmt"
	"github.com/streamingfast/dauth"
	pbauth "github.com/streamingfast/dauth/pb/sf/authentication/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"testing"
	"time"
)

type mockClient struct {
	failOnCount uint64
}

func (m *mockClient) Authenticate(ctx context.Context, in *pbauth.AuthRequest, opts ...grpc.CallOption) (*pbauth.AuthResponse, error) {
	if in.AuthCount == m.failOnCount {
		return nil, fmt.Errorf("authentication failure")
	}
	out := &pbauth.AuthResponse{}
	for _, header := range in.Headers {
		out.AuthenticatedHeaders = append(out.AuthenticatedHeaders, &pbauth.Header{Key: header.Key, Value: header.Value})
	}
	return out, nil

}

func TestAuthenticatorPlugin_ContinuousAuthenticate(t *testing.T) {
	parentCtx := context.Background()
	header := map[string][]string{
		"x-user-id":   []string{"userid"},
		"x-apikey-id": []string{"apiKey"},
	}
	ipAddress := "192.168.1.1"
	continuousInternal := 10 * time.Millisecond
	failOnCount := uint64(3)
	// should be greater than (failOnCount * continuousInternal) + small bugger
	testDuration := 35 * time.Millisecond

	mockClient := &mockClient{failOnCount: failOnCount}
	a := &authenticatorPlugin{
		client:                mockClient,
		continuousInterval:    continuousInternal,
		enabledContinuousAuth: true,
	}

	authenticatedCtx, err := a.Authenticate(parentCtx, "sf.firehose.v1/Blocks", header, ipAddress)
	require.NoError(t, err)

	for {
		select {
		case <-authenticatedCtx.Done():
			require.Equal(t, context.Canceled, authenticatedCtx.Err())
			require.Equal(t, "authentication failure", context.Cause(authenticatedCtx).Error())
			return
		case <-time.After(testDuration):
			assert.Fail(t, "the context should have been canceled by now")
		}
	}

}

type scriptedClient struct {
	responses map[uint64]map[string]string
}

func (m *scriptedClient) Authenticate(ctx context.Context, in *pbauth.AuthRequest, opts ...grpc.CallOption) (*pbauth.AuthResponse, error) {
	out := &pbauth.AuthResponse{}
	for k, v := range m.responses[in.AuthCount] {
		out.AuthenticatedHeaders = append(out.AuthenticatedHeaders, &pbauth.Header{Key: k, Value: v})
	}
	return out, nil
}

func TestAuthenticatorPlugin_ContinuousAuthenticate_RefreshesHeaders(t *testing.T) {
	client := &scriptedClient{responses: map[uint64]map[string]string{
		// other counts get an empty response, like a fail-open authenticator
		2: {"x-user-id": "userid", "x-meta": "meta"},
	}}
	a := &authenticatorPlugin{
		client:                client,
		continuousInterval:    5 * time.Millisecond,
		enabledContinuousAuth: true,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	authenticatedCtx, err := a.Authenticate(ctx, "sf.firehose.v2.Stream/Blocks", nil, "192.168.1.1")
	require.NoError(t, err)
	assert.Equal(t, "", dauth.FromContext(authenticatedCtx).Meta())

	assert.Eventually(t, func() bool {
		return dauth.FromContext(authenticatedCtx).Meta() == "meta"
	}, time.Second, time.Millisecond)

	time.Sleep(25 * time.Millisecond)
	headers := dauth.FromContext(authenticatedCtx)
	assert.Equal(t, "meta", headers.Meta())
	assert.Equal(t, "userid", headers.UserID())
	require.NoError(t, authenticatedCtx.Err())
}
