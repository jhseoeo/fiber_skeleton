package repositoryerror

import "github.com/jhseoeo/fiber-skeleton/src/pkg/typeerr"

var (
	ErrNotFound      = typeerr.NewSentinelError("not found")
	ErrAlreadyExists = typeerr.NewSentinelError("already exists")
)
