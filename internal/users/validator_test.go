package users

import (
	"testing"

	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestValidateSchema(t *testing.T) {
	validator := NewValidator()

	t.Run("Error", func(t *testing.T) {
		assert := assert.New(t)

		user1 := NewUser()
		user2 := NewUser()
		user2.Username = "aaa"
		user2.Name = "a"
		user2.Lastname = "a"
		user3 := NewUser()
		user3.Email = "a"
		user4 := NewUser()
		user4.Username = "admin"
		user4.Password = "admin"
		user4.Email = "a@a"
		user4.Name = "Fulanito"
		user4.Lastname = "De tal"
		user5 := NewUser()
		user5.Username = "admin#$"
		user5.Password = "admin"
		user5.Email = "a$a@%a.com"
		user5.Name = "Fulan~$ito"
		user5.Lastname = "De#tal"
		user6 := NewUser()
		user6.Username = "adm-in"
		user6.Password = "admin"
		user6.Email = "a@a-.com"
		user6.Name = "Fulan1to"
		user6.Lastname = "De tal"
		user7 := NewUser()
		user7.Username = "ádmin"
		user7.Password = "admín"
		user7.Email = "a@-a.com"
		user7.Name = "Fulan1to"
		user7.Lastname = "De0tal"

		tests := []struct {
			in  *User
			out error
		}{{
			in: user1,
			out: errors.Error{
				Type: errors.Validation,
				Code: ErrSchemaValidation.Code,
				Fields: []errors.Field{
					{"username", "invalid", "required"},
					{"password", "invalid", "required"},
					{"email", "invalid", "required"},
					{"name", "invalid", "required"},
					{"lastname", "invalid", "required"},
				},
			},
		}, {
			in: user2,
			out: errors.Error{
				Type: errors.Validation,
				Code: ErrSchemaValidation.Code,
				Fields: []errors.Field{
					{"username", "invalid", "min"},
					{"password", "invalid", "required"},
					{"email", "invalid", "required"},
					{"name", "invalid", "min"},
					{"lastname", "invalid", "min"},
				},
			},
		}, {
			in: user3,
			out: errors.Error{
				Type: errors.Validation,
				Code: ErrSchemaValidation.Code,
				Fields: []errors.Field{
					{"username", "invalid", "required"},
					{"password", "invalid", "required"},
					{"email", "invalid", "email"},
					{"name", "invalid", "required"},
					{"lastname", "invalid", "required"},
				},
			},
		}, {
			in: user4,
			out: errors.Error{
				Type: errors.Validation,
				Code: ErrSchemaValidation.Code,
				Fields: []errors.Field{
					{"email", "invalid", "email"},
				},
			},
		}, {
			in: user5,
			out: errors.Error{
				Type: errors.Validation,
				Code: ErrSchemaValidation.Code,
				Fields: []errors.Field{
					{"username", "invalid", "alphanumdash"},
					{"email", "invalid", "email"},
					{"name", "invalid", "alphaspaces"},
					{"lastname", "invalid", "alphaspaces"},
				},
			},
		}, {
			in: user6,
			out: errors.Error{
				Type: errors.Validation,
				Code: ErrSchemaValidation.Code,
				Fields: []errors.Field{
					{"email", "invalid", "email"},
					{"name", "invalid", "alphaspaces"},
				},
			},
		}, {
			in: user7,
			out: errors.Error{
				Type: errors.Validation,
				Code: ErrSchemaValidation.Code,
				Fields: []errors.Field{
					{"username", "invalid", "alphanumdash"},
					{"email", "invalid", "email"},
					{"name", "invalid", "alphaspaces"},
					{"lastname", "invalid", "alphaspaces"},
				},
			},
		}}

		for _, test := range tests {
			err := validator.ValidateSchema(test.in)
			assert.NotNil(err)
			assert.Equal(err, test.out)
		}
	})

	t.Run("OK", func(t *testing.T) {
		assert := assert.New(t)

		user1 := NewUser()
		user1.Username = "user"
		user1.Password = "pwd"
		user1.Email = "user@email.com"
		user1.Name = "Name"
		user1.Lastname = "Lastname"
		user2 := NewUser()
		user2.Username = "us-er"
		user2.Password = "pwd"
		user2.Email = "user@e-mail.com"
		user2.Name = "Alan Daniel"
		user2.Lastname = "Boglioli Caffé"

		tests := []*User{user1, user2}

		for _, test := range tests {
			err := validator.ValidateSchema(test)
			assert.Nil(err)
		}
	})

}

func TestValidatePassword(t *testing.T) {
	assert := assert.New(t)

	validator := NewValidator()

	tests := []struct {
		in  string
		out error
	}{{
		"123",
		ErrPasswordValidation.F("password", "too_weak"),
	}, {
		"123456",
		ErrPasswordValidation.F("password", "too_weak"),
	}, {
		"1234567",
		ErrPasswordValidation.F("password", "too_weak"),
	}, {
		"12345678",
		nil,
	}}

	for _, test := range tests {
		err := validator.ValidatePassword(test.in)
		assert.Equal(err, test.out)
	}
}
