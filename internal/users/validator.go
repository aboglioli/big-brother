package users

import (
	"regexp"

	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/models"
	govalidator "github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

// Errors
var (
	ErrSchemaValidation   = errors.Validation.New("user.invalid_schema")
	ErrPasswordValidation = errors.Validation.New("user.invalid_password")
)

// Interfaces
type Validator interface {
	ValidateSchema(u *models.User) error
	ValidatePassword(pwd string) error
}

// Implementations
type validator struct {
	validate *govalidator.Validate
}

func NewValidator() Validator {
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

	return &validator{
		validate: validate,
	}
}

func (v *validator) ValidateSchema(u *models.User) error {
	if err := v.validate.Struct(u); err != nil {
		if errs, ok := err.(govalidator.ValidationErrors); ok {
			vErr := ErrSchemaValidation
			for _, err := range errs {
				field := strcase.ToLowerCamel(err.Field())
				vErr = vErr.F(field, "invalid", err.Tag())
			}
			return vErr
		}

		return ErrSchemaValidation
	}

	return nil
}

func (v *validator) ValidatePassword(pwd string) error {
	if len(pwd) < 8 {
		return ErrPasswordValidation.F("password", "too_weak")
	}
	if len(pwd) > 64 {
		return ErrPasswordValidation.F("password", "invalid_length")
	}

	return nil
}
