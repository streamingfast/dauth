package grpc

import (
	"fmt"
	"net/url"
	"time"
)

type config struct {
	endpoint              string
	enabledContinuousAuth bool
	interval              time.Duration
}

func newConfig(urlString string) (*config, error) {
	urlObject, err := url.Parse(urlString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse url: %w", err)
	}

	interval := time.Second * 60
	if intervalParam := urlObject.Query().Get("interval"); intervalParam != "" {
		parsed, err := time.ParseDuration(intervalParam)
		if err != nil {
			return nil, err
		}
		if parsed < time.Second {
			return nil, fmt.Errorf("interval must be at least 1 second")
		}
		interval = parsed
	}

	return &config{
		endpoint:              urlObject.Host + urlObject.Path,
		enabledContinuousAuth: urlObject.Query().Get("continuous") == "true",
		interval:              interval,
	}, nil
}
