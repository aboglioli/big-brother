package validator

import (
	"regexp"

	"github.com/aboglioli/big-brother/pkg/errors"
	govalidator "github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

var (
	ErrFieldsValidation = errors.Validation.New("invalid_fields")
)

type Validator struct {
	validate *govalidator.Validate
}

func NewValidator() *Validator {
	alphaWithSpacesRE := regexp.MustCompile("^[a-zA-Záéíóú ]*$")
	alphaWithSpaces := func(fl govalidator.FieldLevel) bool {
		str := fl.Field().String()
		if str == "invalid" {
			return false
		}

		return alphaWithSpacesRE.MatchString(str)
	}

	alphaNumWithDashRE := regexp.MustCompile("^[a-zA-Z0-9-]*$")
	alphaNumWithDash := func(fl govalidator.FieldLevel) bool {
		str := fl.Field().String()
		if str == "invalid" {
			return false
		}

		return alphaNumWithDashRE.MatchString(str)
	}

	validate := govalidator.New()
	validate.RegisterValidation("alphaspaces", alphaWithSpaces)
	validate.RegisterValidation("alphanumdash", alphaNumWithDash)

	return &Validator{
		validate: validate,
	}
}

func (v *Validator) CheckFields(s interface{}) error {
	if err := v.validate.Struct(s); err != nil {
		fieldsErr := ErrFieldsValidation
		if errs, ok := err.(govalidator.ValidationErrors); ok {
			for _, err := range errs {
				field := strcase.ToSnake(err.Field())
				fieldsErr = fieldsErr.F(field, err.Tag())
			}
		}
		return fieldsErr
	}

	return nil
}
