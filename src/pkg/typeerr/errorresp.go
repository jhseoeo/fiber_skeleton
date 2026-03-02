package typeerr

import (
	"github.com/go-errors/errors"
)

type ErrorResp struct {
	Err     error
	Status  int
	Message string
}

func NewErrorResp(err error, status int, message string) ErrorResp {
	var stackErr *errors.Error
	if !errors.As(err, &stackErr) {
		err = errors.Wrap(err, 1)
	}
	return ErrorResp{
		Err:     err,
		Status:  status,
		Message: message,
	}
}

func (e ErrorResp) Error() string {
	return e.Message
}
