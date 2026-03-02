// custom error code must be added to this file

package errorcode

type ErrorCode int

const (
	Success ErrorCode = 0

	// 400xx - Bad Request
	ErrBadRequest  ErrorCode = 40000
	ErrInvalidID   ErrorCode = 40001
	ErrInvalidBody ErrorCode = 40002

	// 401xx - Unauthorized
	ErrUnauthorized ErrorCode = 40100

	// 404xx - Not Found
	ErrNotFound ErrorCode = 40400

	// 408xx - Request Timeout
	ErrRequestTimeout ErrorCode = 40800

	// 409xx - Conflict
	ErrConflict ErrorCode = 40900

	// 500xx - Internal Server Error
	ErrInternalServer ErrorCode = 50000
)

// HTTPStatus derives the HTTP status code from the error code.
// e.g. 40401 -> 404, 50000 -> 500
func (e ErrorCode) HTTPStatus() int {
	return int(e) / 100
}
