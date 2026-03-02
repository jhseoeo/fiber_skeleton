package validate

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// FieldError represents a single field-level validation failure.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e FieldError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// FieldErrors is a slice of FieldError that also implements the error interface.
type FieldErrors []FieldError

func (fe FieldErrors) Error() string {
	if len(fe) == 0 {
		return ""
	}
	msg := fe[0].Error()
	for _, e := range fe[1:] {
		msg += "; " + e.Error()
	}
	return msg
}

var instance = validator.New()

// Struct validates s and returns nil or a FieldErrors value.
func Struct(s any) error {
	err := instance.Struct(s)
	if err == nil {
		return nil
	}

	var ve validator.ValidationErrors
	if ok := asValidationErrors(err, &ve); !ok {
		return err
	}

	out := make(FieldErrors, 0, len(ve))
	for _, fe := range ve {
		out = append(out, FieldError{
			Field:   fe.Field(),
			Message: fieldMessage(fe),
		})
	}
	return out
}

func asValidationErrors(err error, out *validator.ValidationErrors) bool {
	if ve, ok := err.(validator.ValidationErrors); ok {
		*out = ve
		return true
	}
	return false
}

func fieldMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "min":
		return fmt.Sprintf("must be at least %s", fe.Param())
	case "max":
		return fmt.Sprintf("must be at most %s", fe.Param())
	case "email":
		return "must be a valid email address"
	default:
		return fmt.Sprintf("failed validation: %s", fe.Tag())
	}
}
