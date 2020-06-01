package null

import (
	"github.com/dfuse-io/dauth/ratelimiter"
)

func init() {
	// null://
	ratelimiter.Register("null", func(configURL string) (requestRateLimiter ratelimiter.RateLimiter, err error) {
		return requestRateLimiter, err
	})
}

type RequestRateLimiter struct{}

func NewRequestRateLimiter() *RequestRateLimiter {
	return &RequestRateLimiter{}
}

func (r *RequestRateLimiter) Gate(id string, method string) (allow bool) {
	return true
}
