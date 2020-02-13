package users

import (
	"net/http"

	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/models"
	gValidator "github.com/aboglioli/big-brother/pkg/validator"
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
	validator *gValidator.Validator
}

func NewValidator() Validator {
	return &validator{
		validator: gValidator.NewValidator(),
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
	if err := v.validator.CheckFields(u); err != nil {
		return ErrSchemaValidation.Wrap(err)
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
	if err := v.validator.CheckFields(req); err != nil {
		return errors.ErrRequest.Wrap(err)
	}
	return nil
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

	if err := v.validator.CheckFields(req); err != nil {
		return errors.ErrRequest.Wrap(err)
	}
	return nil
}

func (v *validator) ChangePasswordRequest(req *ChangePasswordRequest) error {
	if err := v.validator.CheckFields(req); err != nil {
		return errors.ErrRequest.Wrap(err)
	}
	return nil
}

func (v *validator) LoginRequest(req *LoginRequest) error {
	if err := v.validator.CheckFields(req); err != nil {
		return errors.ErrRequest.Wrap(err)
	}
	return nil
}
