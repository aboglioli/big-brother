package users

import (
	"net/http"
	"regexp"

	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/models"
	govalidator "github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

// Errors
var (
	ErrUserNotValidated   = errors.Status.New("user.not_validated").S(http.StatusNotFound)
	ErrSchemaValidation   = errors.Validation.New("user.invalid_schema")
	ErrPasswordValidation = errors.Validation.New("user.invalid_password")
)

// Interfaces
type Validator interface {
	Status(u *models.User) error
	Schema(u *models.User) error
	Password(pwd string) error
	RegisterRequest(req *RegisterRequest) error
	UpdateRequest(req *UpdateRequest) error
	ChangePasswordRequest(req *ChangePasswordRequest) error
	LoginRequest(req *LoginRequest) error
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

func (v *validator) Status(u *models.User) error {
	if u == nil {
		return errors.ErrNotFound
	}
	if !u.Enabled {
		return errors.ErrNotFound
	}
	if !u.Validated {
		return ErrUserNotValidated
	}
	return nil
}

func (v *validator) Schema(u *models.User) error {
	if err := v.validate.Struct(u); err != nil {
		valErr := ErrSchemaValidation
		if errs, ok := err.(govalidator.ValidationErrors); ok {
			for _, err := range errs {
				field := strcase.ToSnake(err.Field())
				valErr = valErr.F(field, err.Tag())
			}
			return valErr
		}
		return valErr
	}
	return nil
}

func (v *validator) Password(pwd string) error {
	if len(pwd) < 8 {
		return ErrPasswordValidation.F("password", "too_weak")
	}
	if len(pwd) > 64 {
		return ErrPasswordValidation.F("password", "too_long")
	}

	return nil
}

func (v *validator) RegisterRequest(req *RegisterRequest) error {
	return v.checkFields(req)
}

func (v *validator) UpdateRequest(req *UpdateRequest) error {
	if req == nil ||
		(req.Username == nil &&
			req.Email == nil &&
			req.Name == nil &&
			req.Lastname == nil) {
		return errors.ErrRequest.M("empty request")
	}

	err := errors.ErrRequest
	if req.Username != nil && *req.Username == "" {
		err = err.F("username", "empty")
	}
	if req.Email != nil && *req.Email == "" {
		err = err.F("email", "empty")
	}
	if req.Name != nil && *req.Name == "" {
		err = err.F("name", "empty")
	}
	if req.Lastname != nil && *req.Lastname == "" {
		err = err.F("lastname", "empty")
	}
	if len(err.Fields) > 0 {
		return err
	}

	return v.checkFields(req)
}

func (v *validator) ChangePasswordRequest(req *ChangePasswordRequest) error {
	return v.checkFields(req)
}

func (v *validator) LoginRequest(req *LoginRequest) error {
	return v.checkFields(req)
}

func (v *validator) checkFields(s interface{}) error {
	if err := v.validate.Struct(s); err != nil {
		reqErr := errors.ErrRequest
		if errs, ok := err.(govalidator.ValidationErrors); ok {
			for _, err := range errs {
				field := strcase.ToSnake(err.Field())
				reqErr = reqErr.F(field, err.Tag())
			}
			return reqErr
		}
		return reqErr
	}

	return nil
}
