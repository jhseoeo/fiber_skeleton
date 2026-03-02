package typeerr

import (
	"github.com/go-errors/errors"
	"github.com/jhseoeo/fiber-skeleton/src/dto/errorcode"
)

type ErrorResp struct {
	Err     error
	Code    errorcode.ErrorCode
	Message string
}

func NewErrorResp(err error, code errorcode.ErrorCode, message string) ErrorResp {
	var stackErr *errors.Error
	if !errors.As(err, &stackErr) {
		err = errors.Wrap(err, 1)
	}
	return ErrorResp{
		Err:     err,
		Code:    code,
		Message: message,
	}
}

func (e ErrorResp) Error() string {
	return e.Message
}
