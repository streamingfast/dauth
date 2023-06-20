package grpc

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/streamingfast/dauth"
	pbauth "github.com/streamingfast/dauth/pb/sf/authentication/v1"
	"github.com/streamingfast/dgrpc"
)

func Register() {
	dauth.Register("grpc", func(configURL string) (dauth.Authenticator, error) {
		serverAddr, err := parseURL(configURL)
		if err != nil {
			return nil, fmt.Errorf("grpc factory: %w", err)
		}

		return newAuthenticator(serverAddr)
	})
}

func parseURL(configURL string) (serverAddr string, err error) {
	urlObject, err := url.Parse(configURL)
	if err != nil {
		return
	}
	return urlObject.Host + urlObject.Path, nil
}

type authenticatorPlugin struct {
	client pbauth.AuthenticationClient
}

func newAuthenticator(serverAddr string) (*authenticatorPlugin, error) {
	conn, err := dgrpc.NewInternalNoWaitClient(serverAddr)
	if err != nil {
		return nil, fmt.Errorf("new auth grpc client: %w", err)
	}

	ap := &authenticatorPlugin{
		client: pbauth.NewAuthenticationClient(conn),
	}
	return ap, nil
}

func (a *authenticatorPlugin) Authenticate(ctx context.Context, path string, headers map[string][]string, ipAddress string) (context.Context, error) {
	req := &pbauth.AuthRequest{
		Url:     path,
		Ip:      ipAddress,
		Headers: nil,
	}

	for key, values := range headers {
		for _, value := range values {
			req.Headers = append(req.Headers, &pbauth.Header{
				Key:   strings.ToLower(key),
				Value: value,
			})
		}
	}

	resp, err := a.client.Authenticate(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("auth grpc service failed: %w", err)
	}

	out := make(dauth.TrustedHeaders)
	for _, authenticatedHeader := range resp.AuthenticatedHeaders {
		out[authenticatedHeader.Key] = strings.ToLower(authenticatedHeader.Value)
	}
	return dauth.WithTrustedHeaders(ctx, out), nil
}
