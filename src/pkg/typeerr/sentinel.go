package typeerr

import "github.com/go-errors/errors"

type SentinelError struct {
	error
}

func NewSentinelError(msg string) SentinelError {
	return SentinelError{error: errors.New(msg)}
}

func (e SentinelError) New(innerMsg string) error {
	return errors.Errorf("%w: %s", e, innerMsg)
}
