package grpc

import (
	"context"
	"fmt"
	"testing"
	"testing/synctest"
	"time"

	"github.com/streamingfast/dauth"
	pbauth "github.com/streamingfast/dauth/pb/sf/authentication/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
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
	synctest.Test(t, func(t *testing.T) {
		header := map[string][]string{
			"x-user-id":   {"userid"},
			"x-apikey-id": {"apiKey"},
		}

		a := &authenticatorPlugin{
			client:                &mockClient{failOnCount: 3},
			continuousInterval:    time.Second,
			enabledContinuousAuth: true,
		}

		authenticatedCtx, err := a.Authenticate(context.Background(), "sf.firehose.v1/Blocks", header, "192.168.1.1")
		require.NoError(t, err)

		// AuthCount 2 succeeds, AuthCount 3 fails
		time.Sleep(time.Second)
		synctest.Wait()
		require.NoError(t, authenticatedCtx.Err())

		time.Sleep(time.Second)
		synctest.Wait()
		require.Equal(t, context.Canceled, authenticatedCtx.Err())
		require.Equal(t, "authentication failure", context.Cause(authenticatedCtx).Error())
	})
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
	synctest.Test(t, func(t *testing.T) {
		client := &scriptedClient{responses: map[uint64]map[string]string{
			// other counts get an empty response, like a fail-open authenticator
			2: {"x-user-id": "userid", "x-meta": "meta"},
		}}
		a := &authenticatorPlugin{
			client:                client,
			continuousInterval:    time.Second,
			enabledContinuousAuth: true,
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		authenticatedCtx, err := a.Authenticate(ctx, "sf.firehose.v2.Stream/Blocks", nil, "192.168.1.1")
		require.NoError(t, err)
		assert.Equal(t, "", dauth.FromContext(authenticatedCtx).Meta())

		// AuthCount 2 returns the headers
		time.Sleep(time.Second)
		synctest.Wait()
		assert.Equal(t, "meta", dauth.FromContext(authenticatedCtx).Meta())

		// later empty responses must not remove them
		time.Sleep(5 * time.Second)
		synctest.Wait()
		headers := dauth.FromContext(authenticatedCtx)
		assert.Equal(t, "meta", headers.Meta())
		assert.Equal(t, "userid", headers.UserID())
		require.NoError(t, authenticatedCtx.Err())
	})
}
