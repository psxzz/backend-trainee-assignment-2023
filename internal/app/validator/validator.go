package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type requestValidator struct {
	validator *validator.Validate
}

func New() *requestValidator {
	return &requestValidator{
		validator: validator.New(),
	}
}

func (v *requestValidator) Validate(i any) error {
	if err := v.validator.Struct(i); err != nil {
		return fmt.Errorf("couldn't validate request: %w", err)
	}

	return nil
}
