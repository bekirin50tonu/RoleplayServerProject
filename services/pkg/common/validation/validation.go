package validation

import (
	"github.com/go-playground/validator/v10"
)

type Valitator struct {
	validator *validator.Validate
}

func NewValidator() *Valitator {
	validator := validator.New()
	return &Valitator{
		validator: validator,
	}
}

func (v *Valitator) Validate(data interface{}) error {
	return v.validator.Struct(data)

}
