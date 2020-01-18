package users

import (
	"testing"

	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestValidateSchema(t *testing.T) {
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

func TestValidatePassword(t *testing.T) {
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
