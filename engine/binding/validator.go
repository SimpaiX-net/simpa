package binding

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

/*
Custom validator implementation
*/
type ValidatorImpl interface {
	Engine() any
	ValidateStruct(interface{}) error
}

/*
Default validator
*/
var DefaultValidator ValidatorImpl = &Validator{
	validator.New(),
}

/*
Custom validator implementation (default)
*/
type Validator struct {
	engine *validator.Validate
}

/*
Returns the validator engine
*/
func (v *Validator) Engine() any {
	return v.engine
}

/*
Validates the given struct with the custom engine underlying in v
*/
func (v *Validator) ValidateStruct(data interface{}) error {
	var err error

	err = v.ValidateStruct(data)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return errors.New("InvalidValidationError describes an invalid argument passed to `Struct`, `StructExcept`, StructPartial` or `Field`")
		}

		for _, err = range err.(validator.ValidationErrors) {
			err = errors.New(err.Error())
			break
		}
	}

	return err
}
