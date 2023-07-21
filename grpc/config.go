package grpc

import (
	"fmt"
	"net/url"
)

type config struct {
	endpoint              string
	enabledContinuousAuth bool
}

func newConfig(urlString string) (*config, error) {
	urlObject, err := url.Parse(urlString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse url: %w", err)
	}

	return &config{
		endpoint:              urlObject.Host + urlObject.Path,
		enabledContinuousAuth: urlObject.Query().Get("continuous") == "true",
	}, nil
}
