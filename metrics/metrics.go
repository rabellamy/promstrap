package promstrap

import (
	"github.com/go-playground/validator"
)

// Validate is a helper function to make sure that any required fields exist.
func Validate(opts interface{}) error {
	validate := validator.New()

	return validate.Struct(opts)
}
