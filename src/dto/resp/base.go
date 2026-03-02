package resp

import "github.com/jhseoeo/fiber-skeleton/src/dto/errorcode"

type CommonResp struct {
	Code    errorcode.ErrorCode `json:"code"`
	Message string              `json:"message"`
	Data    any                 `json:"data"`
}
