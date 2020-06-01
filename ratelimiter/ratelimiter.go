package ratelimiter

import (
	"fmt"
	"net/url"
)

type RateLimiter interface {
	Gate(id string, method string) (allow bool)
}

var registry = make(map[string]FactoryFunc)
var servicesRegistry []string

func RegisterService(serviceName string) {
	servicesRegistry = append(servicesRegistry, serviceName)
}

func Validate(serviceName string) bool {
	for _, registeredName := range servicesRegistry {
		if registeredName == serviceName {
			return true
		}
	}
	return false
}

func New(config string) (RateLimiter, error) {
	u, err := url.Parse(config)
	if err != nil {
		return nil, err
	}

	factory := registry[u.Scheme]
	if factory == nil {
		panic(fmt.Sprintf("no ratelimiter plugin named \"%s\" is currently registered", u.Scheme))
	}
	return factory(config)
}

type FactoryFunc func(config string) (RateLimiter, error)

func Register(name string, factory FactoryFunc) {
	registry[name] = factory
}
