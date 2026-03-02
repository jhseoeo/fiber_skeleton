package validate

import "github.com/go-playground/validator/v10"

var instance = validator.New()

func Struct(s any) error {
	return instance.Struct(s)
}
