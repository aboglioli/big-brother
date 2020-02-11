package users

import (
	"testing"

	"github.com/aboglioli/big-brother/internal/auth"
	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/models"
	"github.com/aboglioli/big-brother/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func mockUser() *models.User {
	user := models.NewUser()
	user.Username = "user"
	user.Password = "hashed.password"
	user.Email = "user@user.com"
	user.Name = "Name"
	user.Lastname = "Lastname"
	user.Validated = true
	user.Enabled = true
	return user
}

func TestServiceGetByID(t *testing.T) {
	mUser := mockUser()

	tests := []struct {
		name string
		id   string
		err  error
		mock func(m *mockService)
	}{{
		"empty id",
		"",
		errors.ErrNotFound,
		nil,
	}, {
		"invalid id",
		"123",
		errors.ErrNotFound,
		func(m *mockService) {
			m.repo.On("FindByID", "123").Return(nil, ErrRepositoryNotFound)
		},
	}, {
		"invalid id",
		"abc123",
		errors.ErrNotFound,
		func(m *mockService) {
			m.repo.On("FindByID", "abc123").Return(nil, ErrRepositoryNotFound)
		},
	}, {
		"not found in db",
		mUser.ID,
		errors.ErrNotFound,
		func(m *mockService) {
			m.repo.On("FindByID", mUser.ID).Return(nil, ErrRepositoryNotFound)
		},
	}, {
		"not validated",
		mUser.ID,
		ErrUserNotValidated,
		func(m *mockService) {
			m.repo.On("FindByID", mUser.ID).Return(mUser, nil)
			m.validator.On("Status", mUser).Return(ErrUserNotValidated)
		},
	}, {
		"not enabled",
		mUser.ID,
		errors.ErrNotFound,
		func(m *mockService) {
			m.repo.On("FindByID", mUser.ID).Return(mUser, nil)
			m.validator.On("Status", mUser).Return(errors.ErrNotFound)
		},
	}, {
		"existing user",
		mUser.ID,
		nil,
		func(m *mockService) {
			m.repo.On("FindByID", mUser.ID).Return(mUser, nil)
			m.validator.On("Status", mUser).Return(nil)
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			m := newMockService()
			if test.mock != nil {
				test.mock(m)
			}

			user, err := m.service.GetByID(test.id)

			if test.err != nil { // Error
				errors.Assert(t, test.err, err)
				assert.Nil(user)
			} else { // OK
				assert.Nil(err)
				if assert.NotNil(user) {
					assert.Equal(test.id, user.ID)
					assert.Equal(mUser, user)
				}
			}
			m.crypt.AssertExpectations(t)
			m.repo.AssertExpectations(t)
			m.validator.AssertExpectations(t)
			m.events.AssertExpectations(t)
			m.authServ.AssertExpectations(t)
		})
	}
}

func TestServiceRegister(t *testing.T) {
	mUser := mockUser()
	req := &RegisterRequest{
		Username: mUser.Username,
		Password: "12345678",
		Email:    mUser.Email,
		Name:     mUser.Name,
		Lastname: mUser.Lastname,
	}

	tests := []struct {
		name string
		req  *RegisterRequest
		err  error
		mock func(m *mockService)
	}{{
		"invalid request",
		req,
		errors.ErrRequest,
		func(m *mockService) {
			m.validator.On("RegisterRequest", req).Return(errors.ErrRequest)
		},
	}, {
		"invalid password",
		req,
		ErrPasswordValidation,
		func(m *mockService) {
			m.validator.On("RegisterRequest", req).Return(nil)
			m.validator.On("Password", req.Password).Return(ErrPasswordValidation)
		},
	}, {
		"invalid schema",
		req,
		ErrSchemaValidation,
		func(m *mockService) {
			m.validator.On("RegisterRequest", req).Return(nil)
			m.validator.On("Password", req.Password).Return(nil)
			m.validator.On("Schema", mock.AnythingOfType("*models.User")).Return(ErrSchemaValidation)
		},
	}, {
		"username not available",
		req,
		ErrUserNotAvailable.F("username", "not_available"),
		func(m *mockService) {
			m.validator.On("RegisterRequest", req).Return(nil)
			m.validator.On("Password", req.Password).Return(nil)
			m.validator.On("Schema", mock.AnythingOfType("*models.User")).Return(nil)
			m.repo.On("FindByUsername", req.Username).Return(mUser, nil)
			m.repo.On("FindByEmail", req.Email).Return(nil, ErrRepositoryNotFound)
		},
	}, {
		"email not available",
		req,
		ErrUserNotAvailable.F("email", "not_available"),
		func(m *mockService) {
			m.validator.On("RegisterRequest", req).Return(nil)
			m.validator.On("Password", req.Password).Return(nil)
			m.validator.On("Schema", mock.AnythingOfType("*models.User")).Return(nil)
			m.repo.On("FindByUsername", req.Username).Return(nil, ErrRepositoryNotFound)
			m.repo.On("FindByEmail", req.Email).Return(mUser, nil)
		},
	}, {
		"user and email not available",
		req,
		ErrUserNotAvailable.F("username", "not_available").F("email", "not_available"),
		func(m *mockService) {
			m.validator.On("RegisterRequest", req).Return(nil)
			m.validator.On("Password", req.Password).Return(nil)
			m.validator.On("Schema", mock.AnythingOfType("*models.User")).Return(nil)
			m.repo.On("FindByUsername", req.Username).Return(mockUser(), nil)
			m.repo.On("FindByEmail", req.Email).Return(mockUser(), nil)
		},
	}, {
		"error on insert",
		req,
		errors.ErrInternalServer,
		func(m *mockService) {
			m.validator.On("RegisterRequest", req).Return(nil)
			m.validator.On("Password", req.Password).Return(nil)
			m.validator.On("Schema", mock.AnythingOfType("*models.User")).Return(nil)
			m.repo.On("FindByUsername", req.Username).Return(nil, ErrRepositoryNotFound)
			m.repo.On("FindByEmail", req.Email).Return(nil, ErrRepositoryNotFound)
			m.crypt.On("Hash", req.Password).Return("hashed.password", nil)
			m.repo.On("Insert", mock.AnythingOfType("*models.User")).Return(ErrRepositoryInsert)
		},
	}, {
		"valid user",
		req,
		nil,
		func(m *mockService) {
			m.validator.On("RegisterRequest", req).Return(nil)
			m.validator.On("Password", req.Password).Return(nil)
			m.validator.On("Schema", mock.AnythingOfType("*models.User")).Return(nil)
			m.repo.On("FindByUsername", req.Username).Return(nil, ErrRepositoryNotFound)
			m.repo.On("FindByEmail", req.Email).Return(nil, ErrRepositoryNotFound)
			m.crypt.On("Hash", req.Password).Return("hashed.password", nil)
			m.repo.On("Insert", mock.AnythingOfType("*models.User")).Return(nil)
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			serv := newMockService()
			if test.mock != nil {
				test.mock(serv)
			}

			user, err := serv.Register(test.req)

			if test.err != nil {
				errors.Assert(t, test.err, err)
				assert.Nil(user)
			} else {
				assert.Nil(err)
				if assert.NotNil(user) {
					assert.NotEmpty(user.ID)
					assert.Equal(test.req.Username, user.Username)
					assert.NotEqual(test.req.Password, user.Password)
					assert.Equal("hashed.password", user.Password)
					assert.Equal(test.req.Email, user.Email)
					assert.Equal(test.req.Name, user.Name)
					assert.Equal(test.req.Lastname, user.Lastname)
				}
				serv.validator.AssertCalled(t, "Schema", user)
				serv.repo.AssertCalled(t, "Insert", user)
			}
			serv.crypt.AssertExpectations(t)
			serv.repo.AssertExpectations(t)
			serv.validator.AssertExpectations(t)
			serv.events.AssertExpectations(t)
			serv.authServ.AssertExpectations(t)
		})
	}
}

func TestServiceUpdate(t *testing.T) {
	mUser := mockUser()
	req := &UpdateRequest{
		Username: utils.NewString("new-username"),
		Email:    utils.NewString("new@email.com"),
		Name:     utils.NewString("New name"),
		Lastname: utils.NewString("New lastname"),
	}

	tests := []struct {
		name string
		id   string
		req  *UpdateRequest
		err  error
		mock func(m *mockService)
	}{{
		"invalid request",
		"user123",
		req,
		errors.ErrRequest,
		func(m *mockService) {
			m.validator.On("UpdateRequest", req).Return(errors.ErrRequest)
		},
	}, {
		"not existing user",
		"abc123",
		req,
		errors.ErrNotFound,
		func(m *mockService) {
			m.validator.On("UpdateRequest", req).Return(nil)
			m.repo.On("FindByID", "abc123").Return(nil, ErrRepositoryNotFound)
		},
	}, {
		"not enabled user",
		mUser.ID,
		req,
		errors.ErrNotFound,
		func(m *mockService) {
			m.validator.On("UpdateRequest", req).Return(nil)
			m.repo.On("FindByID", mUser.ID).Return(mUser, nil)
			m.validator.On("Status", mUser).Return(errors.ErrNotFound)
		},
	}, {
		"not validated user",
		mUser.ID,
		req,
		ErrUserNotValidated,
		func(m *mockService) {
			m.validator.On("UpdateRequest", req).Return(nil)
			m.repo.On("FindByID", mUser.ID).Return(mUser, nil)
			m.validator.On("Status", mUser).Return(ErrUserNotValidated)
		},
	}, {
		"invalid name and lastname",
		mUser.ID,
		&UpdateRequest{
			Name:     utils.NewString("New name"),
			Lastname: utils.NewString("New lastname"),
		},
		ErrSchemaValidation,
		func(m *mockService) {
			u := mUser.Clone()
			m.validator.On("UpdateRequest", mock.AnythingOfType("*users.UpdateRequest")).Return(nil)
			m.repo.On("FindByID", mUser.ID).Return(u, nil)
			m.validator.On("Status", u).Return(nil)
			m.validator.On("Schema", u).Return(ErrSchemaValidation)
		},
	}, {
		"invalid username and email",
		mUser.ID,
		&UpdateRequest{
			Username: utils.NewString("new-user"),
			Email:    utils.NewString("new@email.com"),
		},
		ErrSchemaValidation,
		func(m *mockService) {
			u := mUser.Clone()
			m.validator.On("UpdateRequest", mock.AnythingOfType("*users.UpdateRequest")).Return(nil)
			m.repo.On("FindByID", mUser.ID).Return(u, nil)
			m.validator.On("Status", u).Return(nil)
			m.validator.On("Schema", u).Return(ErrSchemaValidation)
		},
	}, {
		"username not available",
		mUser.ID,
		&UpdateRequest{
			Username: utils.NewString("username"),
		},
		ErrUserNotAvailable,
		func(m *mockService) {
			u := mUser.Clone()
			m.validator.On("UpdateRequest", mock.AnythingOfType("*users.UpdateRequest")).Return(nil)
			m.repo.On("FindByID", mUser.ID).Return(u, nil)
			m.validator.On("Status", u).Return(nil)
			m.validator.On("Schema", u).Return(nil)
			m.repo.On("FindByUsername", "username").Return(mockUser(), nil)
		},
	}, {
		"email not available",
		mUser.ID,
		&UpdateRequest{
			Email: utils.NewString("new@email.com"),
		},
		ErrUserNotAvailable,
		func(m *mockService) {
			u := mUser.Clone()
			m.validator.On("UpdateRequest", mock.AnythingOfType("*users.UpdateRequest")).Return(nil)
			m.repo.On("FindByID", mUser.ID).Return(u, nil)
			m.validator.On("Status", u).Return(nil)
			m.validator.On("Schema", u).Return(nil)
			m.repo.On("FindByEmail", "new@email.com").Return(mockUser(), nil)
		},
	}, {
		"username and email not available",
		mUser.ID,
		&UpdateRequest{
			Username: utils.NewString("new-username"),
			Email:    utils.NewString("new@email.com"),
		},
		ErrUserNotAvailable.F("username", "not_available").F("email", "not_available"),
		func(m *mockService) {
			u := mUser.Clone()
			m.validator.On("UpdateRequest", mock.AnythingOfType("*users.UpdateRequest")).Return(nil)
			m.repo.On("FindByID", mUser.ID).Return(u, nil)
			m.validator.On("Status", u).Return(nil)
			m.validator.On("Schema", u).Return(nil)
			m.repo.On("FindByUsername", "new-username").Return(mockUser(), nil)
			m.repo.On("FindByEmail", "new@email.com").Return(mockUser(), nil)
		},
	}, {
		"error on update",
		mUser.ID,
		req,
		errors.ErrInternalServer,
		func(m *mockService) {
			u := mUser.Clone()
			m.validator.On("UpdateRequest", mock.AnythingOfType("*users.UpdateRequest")).Return(nil)
			m.repo.On("FindByID", mUser.ID).Return(u, nil)
			m.validator.On("Status", u).Return(nil)
			m.validator.On("Schema", u).Return(nil)
			m.repo.On("FindByUsername", *req.Username).Return(nil, ErrRepositoryNotFound)
			m.repo.On("FindByEmail", *req.Email).Return(nil, ErrRepositoryNotFound)
			m.repo.On("Update", u).Return(ErrRepositoryUpdate)
		},
	}, {
		"valid update",
		mUser.ID,
		req,
		nil,
		func(m *mockService) {
			u := mUser.Clone()
			m.validator.On("UpdateRequest", mock.AnythingOfType("*users.UpdateRequest")).Return(nil)
			m.repo.On("FindByID", mUser.ID).Return(u, nil)
			m.validator.On("Status", u).Return(nil)
			m.validator.On("Schema", u).Return(nil)
			m.repo.On("FindByUsername", *req.Username).Return(nil, ErrRepositoryNotFound)
			m.repo.On("FindByEmail", *req.Email).Return(nil, ErrRepositoryNotFound)
			m.repo.On("Update", u).Return(nil)
		},
	}, {
		"change name only",
		mUser.ID,
		&UpdateRequest{
			Name: utils.NewString("New name"),
		},
		nil,
		func(m *mockService) {
			u := mUser.Clone()
			m.validator.On("UpdateRequest", mock.AnythingOfType("*users.UpdateRequest")).Return(nil)
			m.repo.On("FindByID", mUser.ID).Return(u, nil)
			m.validator.On("Status", u).Return(nil)
			m.validator.On("Schema", u).Return(nil)
			m.repo.On("Update", u).Return(nil)
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			serv := newMockService()
			if test.mock != nil {
				test.mock(serv)
			}

			user, err := serv.Update(test.id, test.req)

			if test.err != nil {
				errors.Assert(t, test.err, err)
				assert.Nil(user)
			} else {
				assert.Nil(err)
				if assert.NotNil(user) {
					assert.Equal(test.id, user.ID)
					assert.NotEqual(mUser, user)
					if test.req.Username != nil {
						assert.Equal(*test.req.Username, user.Username)
					}
					if test.req.Email != nil {
						assert.Equal(*test.req.Email, user.Email)
					}
					if test.req.Name != nil {
						assert.Equal(*test.req.Name, user.Name)
					}
					if test.req.Lastname != nil {
						assert.Equal(*test.req.Lastname, user.Lastname)
					}
				}
				serv.validator.AssertCalled(t, "UpdateRequest", test.req)
			}
			serv.crypt.AssertExpectations(t)
			serv.repo.AssertExpectations(t)
			serv.validator.AssertExpectations(t)
			serv.events.AssertExpectations(t)
			serv.authServ.AssertExpectations(t)
		})
	}
}

func TestServiceChangePassword(t *testing.T) {
	mUser := mockUser()
	req := &ChangePasswordRequest{
		CurrentPassword: "12345678",
		NewPassword:     "qwertyui",
	}

	tests := []struct {
		name string
		id   string
		req  *ChangePasswordRequest
		err  error
		mock func(m *mockService)
	}{{
		"empty id",
		"",
		req,
		errors.ErrNotFound,
		func(m *mockService) {
			m.validator.On("ChangePasswordRequest", req).Return(nil)
		},
	}, {
		"invalid id",
		"asd123",
		req,
		errors.ErrNotFound,
		func(m *mockService) {
			m.validator.On("ChangePasswordRequest", req).Return(nil)
			m.repo.On("FindByID", "asd123").Return(nil, ErrRepositoryNotFound)
		},
	}, {
		"empty passwords",
		mUser.ID,
		&ChangePasswordRequest{
			CurrentPassword: "",
			NewPassword:     "",
		},
		errors.ErrRequest,
		func(m *mockService) {
			m.validator.On("ChangePasswordRequest", &ChangePasswordRequest{
				CurrentPassword: "",
				NewPassword:     "",
			}).Return(errors.ErrRequest)
		},
	}, {
		"weak new password",
		mUser.ID,
		req,
		ErrPasswordValidation,
		func(m *mockService) {
			m.validator.On("ChangePasswordRequest", req).Return(nil)
			u := mUser.Clone()
			m.repo.On("FindByID", mUser.ID).Return(u, nil)
			m.validator.On("Status", u).Return(nil)
			m.validator.On("Password", req.NewPassword).Return(ErrPasswordValidation)
		},
	}, {
		"mistmatch password",
		mUser.ID,
		req,
		ErrInvalidUser,
		func(m *mockService) {
			m.validator.On("ChangePasswordRequest", req).Return(nil)
			u := mUser.Clone()
			m.repo.On("FindByID", mUser.ID).Return(u, nil)
			m.validator.On("Status", u).Return(nil)
			m.validator.On("Password", req.NewPassword).Return(nil)
			m.crypt.On("Compare", mUser.Password, req.CurrentPassword).Return(false)
		},
	}, {
		"success",
		mUser.ID,
		req,
		nil,
		func(m *mockService) {
			m.validator.On("ChangePasswordRequest", req).Return(nil)
			u := mUser.Clone()
			m.repo.On("FindByID", mUser.ID).Return(u, nil)
			m.validator.On("Status", u).Return(nil)
			m.validator.On("Password", req.NewPassword).Return(nil)
			m.crypt.On("Compare", u.Password, req.CurrentPassword).Return(true)
			m.crypt.On("Hash", req.NewPassword).Return("new.hashed.password", nil)
			m.repo.On("Update", u).Return(nil)
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			serv := newMockService()
			if test.mock != nil {
				test.mock(serv)
			}

			err := serv.ChangePassword(test.id, test.req)

			errors.Assert(t, test.err, err)

			serv.validator.AssertExpectations(t)
			serv.repo.AssertExpectations(t)
			serv.crypt.AssertExpectations(t)
		})
	}
}

func TestServiceDelete(t *testing.T) {
	mUser := mockUser()

	tests := []struct {
		name string
		id   string
		err  error
		mock func(m *mockService)
	}{{
		"not found",
		"abc123",
		errors.ErrNotFound,
		func(m *mockService) {
			m.repo.On("FindByID", "abc123").Return(nil, ErrRepositoryNotFound)
		},
	}, {
		"not enabled",
		mUser.ID,
		errors.ErrNotFound,
		func(m *mockService) {
			u := mUser.Clone()
			u.Enabled = false
			m.repo.On("FindByID", mUser.ID).Return(u, nil)
			m.validator.On("Status", u).Return(errors.ErrNotFound)
		},
	}, {
		"not validated",
		mUser.ID,
		ErrUserNotValidated,
		func(m *mockService) {
			u := mUser.Clone()
			u.Validated = false
			m.repo.On("FindByID", mUser.ID).Return(u, nil)
			m.validator.On("Status", u).Return(ErrUserNotValidated)
		},
	}, {
		"error on delete",
		mUser.ID,
		errors.ErrInternalServer,
		func(m *mockService) {
			u := mUser.Clone()
			m.repo.On("FindByID", mUser.ID).Return(u, nil)
			m.validator.On("Status", u).Return(nil)
			m.repo.On("Delete", mUser.ID).Return(ErrRepositoryDelete)
		},
	}, {
		"success",
		mUser.ID,
		nil,
		func(m *mockService) {
			u := mUser.Clone()
			m.repo.On("FindByID", mUser.ID).Return(u, nil)
			m.validator.On("Status", u).Return(nil)
			m.repo.On("Delete", mUser.ID).Return(nil)
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			serv := newMockService()
			if test.mock != nil {
				test.mock(serv)
			}

			err := serv.Delete(test.id)

			if test.err != nil {
				if assert.NotNil(err) {
					errors.Assert(t, test.err, err)
				}
			} else {
				assert.Nil(err)
			}
			serv.crypt.AssertExpectations(t)
			serv.repo.AssertExpectations(t)
			serv.validator.AssertExpectations(t)
			serv.events.AssertExpectations(t)
			serv.authServ.AssertExpectations(t)
		})
	}
}

func TestServiceLogin(t *testing.T) {
	mUser := mockUser()
	mTokenStr := "encoded.token"

	req := &LoginRequest{
		UsernameOrEmail: "user",
		Password:        "12345678",
	}

	tests := []struct {
		name string
		req  *LoginRequest
		err  error
		mock func(m *mockService)
	}{{
		"empty request",
		req,
		errors.ErrRequest,
		func(m *mockService) {
			m.validator.On("LoginRequest", req).Return(errors.ErrRequest)
		},
	}, {
		"invalid username",
		req,
		ErrInvalidUser,
		func(m *mockService) {
			m.validator.On("LoginRequest", req).Return(nil)
			m.repo.On("FindByUsername", req.UsernameOrEmail).Return(nil, ErrRepositoryNotFound)
			m.repo.On("FindByEmail", req.UsernameOrEmail).Return(nil, ErrRepositoryNotFound)
		},
	}, {
		"invalid password",
		req,
		ErrInvalidUser,
		func(m *mockService) {
			m.validator.On("LoginRequest", req).Return(nil)
			m.repo.On("FindByUsername", req.UsernameOrEmail).Return(nil, ErrRepositoryNotFound)
			m.repo.On("FindByEmail", req.UsernameOrEmail).Return(mUser, nil)
			m.crypt.On("Compare", mUser.Password, req.Password).Return(false)
		},
	}, {
		"login with username and password",
		req,
		nil,
		func(m *mockService) {
			m.validator.On("LoginRequest", req).Return(nil)
			m.repo.On("FindByUsername", req.UsernameOrEmail).Return(mUser, nil)
			m.crypt.On("Compare", mUser.Password, req.Password).Return(true)
			m.authServ.On("Create", mUser.ID).Return(mTokenStr, nil)
		},
	}, {
		"login with email and password",
		req,
		nil,
		func(m *mockService) {
			m.validator.On("LoginRequest", req).Return(nil)
			m.repo.On("FindByUsername", req.UsernameOrEmail).Return(nil, ErrRepositoryNotFound)
			m.repo.On("FindByEmail", req.UsernameOrEmail).Return(mUser, nil)
			m.crypt.On("Compare", mUser.Password, req.Password).Return(true)
			m.authServ.On("Create", mUser.ID).Return(mTokenStr, nil)
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			serv := newMockService()
			if test.mock != nil {
				test.mock(serv)
			}

			tokenStr, err := serv.Login(test.req)

			if test.err != nil {
				errors.Assert(t, test.err, err)
				assert.Empty(tokenStr)
			} else {
				assert.Nil(err)
				if assert.NotEmpty(tokenStr) {
					assert.Equal(mTokenStr, tokenStr)
				}
			}
			serv.crypt.AssertExpectations(t)
			serv.repo.AssertExpectations(t)
			serv.validator.AssertExpectations(t)
			serv.events.AssertExpectations(t)
			serv.authServ.AssertExpectations(t)
		})
	}
}

func TestServiceLogout(t *testing.T) {
	mUser := mockUser()
	mToken := models.NewToken(mUser.ID)
	mTokenStr := "encoded.token"

	tests := []struct {
		name     string
		tokenStr string
		err      error
		mock     func(m *mockService)
	}{{
		"empty token",
		"",
		ErrInvalidUser,
		func(m *mockService) {
			m.authServ.On("Invalidate", "").Return(nil, auth.ErrInvalidate)
		},
	}, {
		"non existing user",
		mTokenStr,
		ErrInvalidUser,
		func(m *mockService) {
			m.authServ.On("Invalidate", mTokenStr).Return(mToken, nil)
			m.repo.On("FindByID", mUser.ID).Return(nil, ErrRepositoryNotFound)
		},
	}, {
		"success",
		mTokenStr,
		nil,
		func(m *mockService) {
			m.authServ.On("Invalidate", mTokenStr).Return(mToken, nil)
			m.repo.On("FindByID", mUser.ID).Return(mUser, nil)
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			serv := newMockService()
			if test.mock != nil {
				test.mock(serv)
			}

			err := serv.Logout(test.tokenStr)

			errors.Assert(t, test.err, err)

			serv.repo.AssertExpectations(t)
			serv.authServ.AssertExpectations(t)
			serv.events.AssertExpectations(t)
			serv.validator.AssertExpectations(t)
			serv.crypt.AssertExpectations(t)
		})
	}
}
