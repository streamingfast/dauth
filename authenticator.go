package dauth

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/url"
)

type Authenticator interface {
	Authenticate(ctx context.Context, path string, headers map[string][]string, ipAddress string) (context.Context, error)
	Ready(context.Context) bool
}

var registry = make(map[string]FactoryFunc)

func New(config string, logger *zap.Logger) (Authenticator, error) {
	u, err := url.Parse(config)
	if err != nil {
		return nil, err
	}

	factory := registry[u.Scheme]
	if factory == nil {
		panic(fmt.Sprintf("no Authenticator plugin named \"%s\" is currently registered", u.Scheme))
	}
	return factory(config, logger)
}

type FactoryFunc func(config string, logger *zap.Logger) (Authenticator, error)

func Register(name string, factory FactoryFunc) {
	registry[name] = factory
}
