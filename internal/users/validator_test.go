package users

import (
	"testing"

	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/models"
	"github.com/aboglioli/big-brother/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestValidatorStatus(t *testing.T) {
	tests := []struct {
		name string
		user func(u *models.User) *models.User
		err  error
	}{{
		"nil user",
		func(u *models.User) *models.User {
			return nil
		},
		errors.ErrNotFound,
	}, {
		"not enabled user",
		func(u *models.User) *models.User {
			u.Validated = true
			u.Enabled = false
			return u
		},
		errors.ErrNotFound,
	}, {
		"not validated user",
		func(u *models.User) *models.User {
			u.Validated = false
			u.Enabled = true
			return u
		},
		ErrUserNotValidated,
	}, {
		"enabled and validated user",
		func(u *models.User) *models.User {
			u.Validated = true
			u.Enabled = true
			return u
		},
		nil,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mUser := models.NewUser()
			if test.user != nil {
				mUser = test.user(mUser)
			}
			validator := NewValidator()

			err := validator.Status(mUser)

			errors.Assert(t, test.err, err)
		})
	}
}

func TestValidatorSchema(t *testing.T) {
	tests := []struct {
		name string
		mock func(u *models.User)
		err  error
	}{{
		"empty user",
		func(u *models.User) {
			eu := models.NewUser()
			*u = *eu
		},
		ErrSchemaValidation,
	}, {
		"empty username",
		func(u *models.User) { u.Username = "" },
		ErrSchemaValidation,
	}, {
		"invalid username",
		func(u *models.User) { u.Username = "aaa" },
		ErrSchemaValidation,
	}, {
		"invalid username",
		func(u *models.User) { u.Username = "admin#$" },
		ErrSchemaValidation,
	}, {
		"invalid username",
		func(u *models.User) { u.Username = "@dmin" },
		ErrSchemaValidation,
	}, {
		"invalid username",
		func(u *models.User) { u.Username = "ádmin" },
		ErrSchemaValidation,
	}, {
		"empty password",
		func(u *models.User) { u.Password = "" },
		ErrSchemaValidation,
	}, {
		"empty email",
		func(u *models.User) { u.Email = "" },
		ErrSchemaValidation,
	}, {
		"invalid email",
		func(u *models.User) { u.Email = "a" },
		ErrSchemaValidation,
	}, {
		"invalid email",
		func(u *models.User) { u.Email = "a@a" },
		ErrSchemaValidation,
	}, {
		"invalid email",
		func(u *models.User) { u.Email = "a@a-.com" },
		ErrSchemaValidation,
	}, {
		"invalid email",
		func(u *models.User) { u.Email = "a@-a.com" },
		ErrSchemaValidation,
	}, {
		"invalid name",
		func(u *models.User) { u.Name = "Fulan1to" },
		ErrSchemaValidation,
	}, {
		"invalid name",
		func(u *models.User) { u.Name = "Ful@ano" },
		ErrSchemaValidation,
	}, {
		"invalid lastname",
		func(u *models.User) { u.Lastname = "De t@l" },
		ErrSchemaValidation,
	}, {
		"invalid lastname",
		func(u *models.User) { u.Lastname = "De0tal" },
		ErrSchemaValidation,
	}, {
		"valid",
		nil,
		nil,
	}, {
		"valid",
		func(u *models.User) {
			u.Username = "user-name"
			u.Password = "pwd"
			u.Email = "user@e-mail.com"
			u.Name = "Fulano"
			u.Lastname = "De tal"
		},
		nil,
	}, {
		"accent mark",
		func(u *models.User) {
			u.Name = "Alán"
			u.Lastname = "Boglioli Caffé"
		},
		nil,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			user := mockUser()
			if test.mock != nil {
				test.mock(user)
			}
			validator := NewValidator()
			err := validator.Schema(user)

			if test.err != nil { // Error
				if assert.NotNil(err) {
					errors.Assert(t, test.err, err)
				}
			} else { // OK
				assert.Nil(err)
			}
		})
	}
}

func TestValidatorPassword(t *testing.T) {
	tests := []struct {
		pwd string
		err error
	}{{
		"123",
		ErrPasswordValidation,
	}, {
		"123456",
		ErrPasswordValidation,
	}, {
		"abc123",
		ErrPasswordValidation,
	}, {
		"12345678",
		nil,
	}, {
		"123456789",
		nil,
	}, {
		"long-password#!",
		nil,
	}, {
		"my-compl€x_p@ssw0rd!",
		nil,
	}}

	for _, test := range tests {
		t.Run(test.pwd, func(t *testing.T) {
			assert := assert.New(t)
			validator := NewValidator()
			err := validator.Password(test.pwd)

			if test.err != nil { // Error
				if assert.NotNil(err) {
					errors.Assert(t, test.err, err)
				}
			} else { // OK
				assert.Nil(err)
			}
		})
	}
}

func TestValidatorRegisterRequest(t *testing.T) {
	tests := []struct {
		name string
		req  *RegisterRequest
		err  error
	}{{
		"nil request",
		nil,
		errors.ErrRequest,
	}, {
		"empty request",
		&RegisterRequest{},
		errors.ErrRequest,
	}, {
		"empty name and lastname",
		&RegisterRequest{
			Username: "username",
			Password: "12345678",
			Email:    "user@user.com",
		},
		errors.ErrRequest.F("name", "required").F("lastname", "required"),
	}, {
		"invalid email",
		&RegisterRequest{
			Username: "username",
			Password: "12345678",
			Email:    "á@-a.c",
			Name:     "Name",
			Lastname: "Lastname",
		},
		errors.ErrRequest.F("email", "email"),
	}, {
		"empty password",
		&RegisterRequest{
			Username: "username",
			Password: "",
			Email:    "user@user.com",
			Name:     "Name",
			Lastname: "Lastname",
		},
		errors.ErrRequest.F("password", "required"),
	}, {
		"valid",
		&RegisterRequest{
			Username: "username",
			Password: "12345678",
			Email:    "user@user.com",
			Name:     "Name",
			Lastname: "Lastname",
		},
		nil,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validator := NewValidator()

			err := validator.RegisterRequest(test.req)

			errors.Assert(t, test.err, err)
		})
	}
}

func TestValidatortUpdateRequest(t *testing.T) {
	tests := []struct {
		name string
		req  *UpdateRequest
		err  error
	}{{
		"nil request",
		nil,
		errors.ErrRequest,
	}, {
		"empty request",
		&UpdateRequest{
			Username: utils.NewString(""),
			Email:    utils.NewString(""),
			Name:     utils.NewString(""),
			Lastname: utils.NewString(""),
		},
		errors.ErrRequest,
	}, {
		"empty name and lastname",
		&UpdateRequest{
			Name:     utils.NewString(""),
			Lastname: utils.NewString(""),
		},
		errors.ErrRequest,
	}, {
		"empty username",
		&UpdateRequest{
			Username: utils.NewString(""),
		},
		errors.ErrRequest,
	}, {
		"invalid email",
		&UpdateRequest{
			Email: utils.NewString("à@a-.b"),
		},
		errors.ErrRequest,
	}, {
		"valid",
		&UpdateRequest{
			Username: utils.NewString("username"),
			Email:    utils.NewString("user@user.com"),
			Name:     utils.NewString("Name"),
			Lastname: utils.NewString("Lastname"),
		},
		nil,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validator := NewValidator()

			err := validator.UpdateRequest(test.req)

			errors.Assert(t, test.err, err)
		})
	}
}

func TestValidatorChangePasswordRequest(t *testing.T) {
	tests := []struct {
		name string
		req  *ChangePasswordRequest
		err  error
	}{{
		"nil request",
		nil,
		errors.ErrRequest,
	}, {
		"empty request",
		&ChangePasswordRequest{},
		errors.ErrRequest,
	}, {
		"only current password",
		&ChangePasswordRequest{
			CurrentPassword: "1234",
		},
		errors.ErrRequest,
	}, {
		"only new password",
		&ChangePasswordRequest{
			NewPassword: "abcd",
		},
		errors.ErrRequest,
	}, {
		"valid",
		&ChangePasswordRequest{
			CurrentPassword: "1234",
			NewPassword:     "abcd",
		},
		nil,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validator := NewValidator()

			err := validator.ChangePasswordRequest(test.req)

			errors.Assert(t, test.err, err)
		})
	}
}

func TestValidatorLoginRequest(t *testing.T) {
	tests := []struct {
		name string
		req  *LoginRequest
		err  error
	}{{
		"nil request",
		nil,
		errors.ErrRequest,
	}, {
		"empty request",
		&LoginRequest{},
		errors.ErrRequest,
	}, {
		"only username",
		&LoginRequest{
			UsernameOrEmail: "username",
		},
		errors.ErrRequest,
	}, {
		"only password",
		&LoginRequest{
			Password: "1234",
		},
		errors.ErrRequest,
	}, {
		"valid",
		&LoginRequest{
			UsernameOrEmail: "username",
			Password:        "1234",
		},
		nil,
	}}

	for _, test := range tests {
		validator := NewValidator()

		err := validator.LoginRequest(test.req)

		errors.Assert(t, test.err, err)
	}
}
