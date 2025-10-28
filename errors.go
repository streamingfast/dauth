package dauth

import "errors"

type ErrInvalidAuthentication struct {
	inner error
}

func (e *ErrInvalidAuthentication) Error() string {
	return e.inner.Error()
}

func (e *ErrInvalidAuthentication) Unwrap() error {
	return e.inner
}

func NewErrInvalidAuthentication(inner error) error {
	if inner == nil {
		return nil
	}

	return &ErrInvalidAuthentication{inner: inner}
}

var (
	ErrAuthorizationHeaderNotFound      = NewErrInvalidAuthentication(errors.New("no authorization header found"))
	ErrAuthorizationHeaderMissingBearer = NewErrInvalidAuthentication(errors.New("authorization header format must be Bearer {token}"))
)

// IsErrInvalidAuthentication checks whether the provided error is of type ErrInvalidAuthentication.
func IsErrInvalidAuthentication(err error) bool {
	var target *ErrInvalidAuthentication
	return errors.As(err, &target)
}
