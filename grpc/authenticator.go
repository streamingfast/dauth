package grpc

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/streamingfast/dauth"
	pbauth "github.com/streamingfast/dauth/pb/sf/authentication/v1"
	"github.com/streamingfast/dgrpc"
	pbhealth "google.golang.org/grpc/health/grpc_health_v1"
)

func Register() {
	dauth.Register("grpc", func(configURL string) (dauth.Authenticator, error) {
		config, err := newConfig(configURL)
		if err != nil {
			return nil, fmt.Errorf("failed to setup config: %w", err)
		}
		return newAuthenticator(config)
	})
}

type authenticatorPlugin struct {
	client                pbauth.AuthenticationClient
	healthClient          pbhealth.HealthClient
	continuousInterval    time.Duration
	enabledContinuousAuth bool
}

func newAuthenticator(c *config) (*authenticatorPlugin, error) {
	conn, err := dgrpc.NewInternalNoWaitClient(c.endpoint)
	if err != nil {
		return nil, fmt.Errorf("new auth grpc client: %w", err)
	}

	ap := &authenticatorPlugin{
		client:                pbauth.NewAuthenticationClient(conn),
		enabledContinuousAuth: c.enabledContinuousAuth,
		continuousInterval:    10 * time.Second,
		healthClient:          pbhealth.NewHealthClient(conn),
	}
	return ap, nil
}

func (a *authenticatorPlugin) Ready(ctx context.Context) bool {
	r, err := a.healthClient.Check(ctx, &pbhealth.HealthCheckRequest{})
	if err != nil {
		return false
	}
	return r.Status == pbhealth.HealthCheckResponse_SERVING
}

func (a *authenticatorPlugin) Authenticate(ctx context.Context, path string, headers map[string][]string, ipAddress string) (context.Context, error) {
	req := &pbauth.AuthRequest{
		Url:       path,
		Ip:        ipAddress,
		AuthCount: 1,
		Headers:   nil,
	}

	for key, values := range headers {
		for _, value := range values {
			if !utf8.ValidString(key) {
				key = strings.ToValidUTF8(key, "?")
			}
			if !utf8.ValidString(value) {
				value = strings.ToValidUTF8(value, "?")
			}

			req.Headers = append(req.Headers, &pbauth.Header{
				Key:   strings.ToLower(key),
				Value: value,
			})
		}
	}

	resp, err := a.client.Authenticate(ctx, req)
	if err != nil {
		return nil, err
	}

	out := make(dauth.TrustedHeaders)
	for _, authenticatedHeader := range resp.AuthenticatedHeaders {
		out[strings.ToLower(authenticatedHeader.Key)] = authenticatedHeader.Value
	}

	ctx, cancel := context.WithCancelCause(ctx)

	go a.continuousAuth(ctx, req, cancel)

	return dauth.WithTrustedHeaders(ctx, out), nil
}

func (a *authenticatorPlugin) continuousAuth(ctx context.Context, req *pbauth.AuthRequest, cancel context.CancelCauseFunc) {
	if !a.enabledContinuousAuth {
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(a.continuousInterval):
		}

		req.AuthCount++

		if _, err := a.client.Authenticate(context.Background(), req); err != nil {
			cancel(err)
			return
		}
	}
}
