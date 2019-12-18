package users

import (
	"regexp"

	"github.com/aboglioli/big-brother/errors"
	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

var (
	ErrSchemaValidation   = errors.Validation.New("user.invalid_schema")
	ErrPasswordValidation = errors.Validation.New("user.invalid_password")
)

var ValidateSchema = func(u *User) error {
	validate := validator.New()
	validate.RegisterValidation("alphaspaces", alphaWithSpaces)
	validate.RegisterValidation("alphanumdash", alphaNumWithDash)

	if err := validate.Struct(u); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
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

var ValidatePassword = func(pwd string) error {
	if len(pwd) < 8 {
		return ErrPasswordValidation.F("password", "too_weak")
	}
	if len(pwd) > 64 {
		return ErrPasswordValidation.F("password", "invalid_length")
	}

	return nil
}

var alphaWithSpacesRE = regexp.MustCompile("^[a-zA-Záéíóú ]*$")

func alphaWithSpaces(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	if str == "invalid" {
		return false
	}

	return alphaWithSpacesRE.MatchString(str)
}

var alphaNumWithDashRE = regexp.MustCompile("^[a-zA-Z0-9-]*$")

func alphaNumWithDash(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	if str == "invalid" {
		return false
	}

	return alphaNumWithDashRE.MatchString(str)
}
