package users

import (
	"regexp"

	"github.com/aboglioli/big-brother/errors"
	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

var (
	ErrSchemaValidation = errors.Validation.New("user.schema")
)

func ValidateSchema(u *User) error {
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

		return errors.Unknown.New("user.schema")
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
