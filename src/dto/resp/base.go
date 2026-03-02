package resp

import "github.com/jhseoeo/fiber-skeleton/src/dto/errorcode"

type CommonResp struct {
	Code    errorcode.ErrorCode `json:"code"`
	Message string              `json:"message"`
	Data    any                 `json:"data"`
}

type PaginatedResp struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Data  any `json:"data"`
}
