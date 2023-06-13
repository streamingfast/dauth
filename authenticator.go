package dauth

import (
	"context"
	"fmt"
	"net/url"
)

type Authenticator interface {
	Authenticate(ctx context.Context, path string, headers url.Values, ipAddress string) (url.Values, error)
}

var registry = make(map[string]FactoryFunc)

func New(config string) (Authenticator, error) {
	u, err := url.Parse(config)
	if err != nil {
		return nil, err
	}

	factory := registry[u.Scheme]
	if factory == nil {
		panic(fmt.Sprintf("no Authenticator plugin named \"%s\" is currently registered", u.Scheme))
	}
	return factory(config)
}

type FactoryFunc func(config string) (Authenticator, error)

func Register(name string, factory FactoryFunc) {
	registry[name] = factory
}
