package typeerr

import (
	"github.com/go-errors/errors"
	"github.com/jhseoeo/fiber-skeleton/src/dto/errorcode"
)

type ErrorResp struct {
	Err     error
	Code    errorcode.ErrorCode
	Message string
	Data    any
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

// NewErrorRespWithData is like NewErrorResp but attaches structured data to the response body.
func NewErrorRespWithData(err error, code errorcode.ErrorCode, message string, data any) ErrorResp {
	r := NewErrorResp(err, code, message)
	r.Data = data
	return r
}

func (e ErrorResp) Error() string {
	return e.Message
}
