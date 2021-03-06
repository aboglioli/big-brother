package users

import (
	"testing"

	"github.com/aboglioli/big-brother/internal/auth"
	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/events"
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

func copyUser(u *models.User) *models.User {
	copy := *u
	return &copy
}

func TestGetByID(t *testing.T) {
	mUser := mockUser()

	tests := []struct {
		name string
		id   string
		err  error
		mock func(s *mockService)
	}{{
		"invalid id",
		"123",
		ErrNotFound,
		func(s *mockService) {
			s.repo.On("FindByID", "123").Return(nil, ErrRepositoryNotFound)
		},
	}, {
		"invalid id",
		"abc123",
		ErrNotFound,
		func(s *mockService) {
			s.repo.On("FindByID", "abc123").Return(nil, ErrRepositoryNotFound)
		},
	}, {
		"not found in db",
		mUser.ID,
		ErrNotFound.Wrap(ErrRepositoryNotFound),
		func(s *mockService) {
			s.repo.On("FindByID", mUser.ID).Return(nil, ErrRepositoryNotFound)
		},
	}, {
		"not validated",
		mUser.ID,
		ErrNotValidated,
		func(s *mockService) {
			u := copyUser(mUser)
			u.Validated = false
			s.repo.On("FindByID", mUser.ID).Return(u, nil)
		},
	}, {
		"not enabled",
		mUser.ID,
		ErrNotFound,
		func(s *mockService) {
			u := copyUser(mUser)
			u.Enabled = false
			s.repo.On("FindByID", mUser.ID).Return(u, nil)
		},
	}, {
		"existing user",
		mUser.ID,
		nil,
		func(s *mockService) {
			s.repo.On("FindByID", mUser.ID).Return(mUser, nil)
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			serv := newMockService()
			if test.mock != nil {
				test.mock(serv)
			}

			user, err := serv.GetByID(test.id)

			if test.err != nil { // Error
				if assert.NotNil(err) {
					errors.Assert(t, test.err, err)
				}
				assert.Nil(user)
			} else { // OK
				assert.Nil(err)
				if assert.NotNil(user) {
					assert.Equal(test.id, user.ID)
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

func TestRegister(t *testing.T) {
	mUser := mockUser()

	genReq := func(cb func(req *RegisterRequest)) *RegisterRequest {
		req := &RegisterRequest{
			Username: mUser.Username,
			Password: "12345678",
			Email:    mUser.Email,
			Name:     mUser.Name,
			Lastname: mUser.Lastname,
		}
		if cb != nil {
			cb(req)
		}
		return req
	}

	tests := []struct {
		name string
		req  *RegisterRequest
		err  error
		mock func(s *mockService)
	}{{
		"invalid schema and password",
		genReq(nil),
		errors.Errors{ErrPasswordValidation, ErrSchemaValidation},
		func(s *mockService) {
			s.validator.On("ValidateSchema", mock.AnythingOfType("*models.User")).Return(ErrSchemaValidation)
			s.validator.On("ValidatePassword", "12345678").Return(ErrPasswordValidation)
		},
	}, {
		"invalid password",
		genReq(func(req *RegisterRequest) {
			req.Password = "1234567"
		}),
		errors.Errors{ErrPasswordValidation},
		func(s *mockService) {
			s.validator.On("ValidateSchema", mock.AnythingOfType("*models.User")).Return(nil)
			s.validator.On("ValidatePassword", "1234567").Return(ErrPasswordValidation)
		},
	}, {
		"username not available",
		genReq(nil),
		ErrNotAvailable.F("username", "not_available"),
		func(s *mockService) {
			s.repo.On("FindByUsername", "user").Return(mUser, nil)
			s.repo.On("FindByEmail", "user@user.com").Return(nil, ErrRepositoryNotFound)
			s.validator.On("ValidateSchema", mock.AnythingOfType("*models.User")).Return(nil)
			s.validator.On("ValidatePassword", "12345678").Return(nil)
		},
	}, {
		"email not available",
		genReq(nil),
		ErrNotAvailable.F("email", "not_available"),
		func(s *mockService) {
			s.repo.On("FindByUsername", "user").Return(nil, ErrRepositoryNotFound)
			s.repo.On("FindByEmail", "user@user.com").Return(mUser, nil)
			s.validator.On("ValidateSchema", mock.AnythingOfType("*models.User")).Return(nil)
			s.validator.On("ValidatePassword", "12345678").Return(nil)
		},
	}, {
		"user and email not available",
		genReq(nil),
		ErrNotAvailable.F("username", "not_available").F("email", "not_available"),
		func(s *mockService) {
			s.repo.On("FindByUsername", "user").Return(mUser, nil)
			s.repo.On("FindByEmail", "user@user.com").Return(mUser, nil)
			s.validator.On("ValidateSchema", mock.AnythingOfType("*models.User")).Return(nil)
			s.validator.On("ValidatePassword", "12345678").Return(nil)
		},
	}, {
		"error on insert",
		genReq(nil),
		ErrRegister.Wrap(ErrRepositoryInsert),
		func(s *mockService) {
			s.repo.On("FindByUsername", "user").Return(nil, ErrRepositoryNotFound)
			s.repo.On("FindByEmail", "user@user.com").Return(nil, ErrRepositoryNotFound)
			s.validator.On("ValidateSchema", mock.AnythingOfType("*models.User")).Return(nil)
			s.validator.On("ValidatePassword", "12345678").Return(nil)
			s.crypt.On("Hash", "12345678").Return("hashed.password", nil)
			s.repo.On("Insert", mock.AnythingOfType("*models.User")).Return(ErrRepositoryInsert)
		},
	}, {
		"error on publishing event",
		genReq(nil),
		ErrRegister.Wrap(events.ErrPublish),
		func(s *mockService) {
			s.repo.On("FindByUsername", "user").Return(nil, ErrRepositoryNotFound)
			s.repo.On("FindByEmail", "user@user.com").Return(nil, ErrRepositoryNotFound)
			s.validator.On("ValidateSchema", mock.AnythingOfType("*models.User")).Return(nil)
			s.validator.On("ValidatePassword", "12345678").Return(nil)
			s.crypt.On("Hash", "12345678").Return("hashed.password", nil)
			s.repo.On("Insert", mock.AnythingOfType("*models.User")).Return(nil)
			s.events.On("Publish", mock.AnythingOfType("*users.UserEvent"), mock.AnythingOfType("*events.Options")).Return(events.ErrPublish)
		},
	}, {
		"valid user",
		genReq(nil),
		nil,
		func(s *mockService) {
			s.repo.On("FindByUsername", "user").Return(nil, ErrRepositoryNotFound)
			s.repo.On("FindByEmail", "user@user.com").Return(nil, ErrRepositoryNotFound)
			s.validator.On("ValidateSchema", mock.AnythingOfType("*models.User")).Return(nil)
			s.validator.On("ValidatePassword", "12345678").Return(nil)
			s.crypt.On("Hash", "12345678").Return("hashed.password", nil)
			s.repo.On("Insert", mock.AnythingOfType("*models.User")).Return(nil)
			s.events.On("Publish", mock.AnythingOfType("*users.UserEvent"), mock.AnythingOfType("*events.Options")).Return(nil)
		},
	}, {
		"valid admin",
		genReq(func(req *RegisterRequest) {
			req.Username = "admin"
			req.Password = "adminComplexPasswd#!"
			req.Email = "admin@admin.com"
			req.Name = "Admin"
			req.Lastname = "Lastname"
		}),
		nil,
		func(s *mockService) {
			s.repo.On("FindByUsername", "admin").Return(nil, ErrRepositoryNotFound)
			s.repo.On("FindByEmail", "admin@admin.com").Return(nil, ErrRepositoryNotFound)
			s.validator.On("ValidateSchema", mock.AnythingOfType("*models.User")).Return(nil)
			s.validator.On("ValidatePassword", "adminComplexPasswd#!").Return(nil)
			s.crypt.On("Hash", "adminComplexPasswd#!").Return("hashed.password", nil)
			s.repo.On("Insert", mock.AnythingOfType("*models.User")).Return(nil)
			s.events.On("Publish", mock.AnythingOfType("*users.UserEvent"), mock.AnythingOfType("*events.Options")).Return(nil)
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
				if assert.NotNil(err) {
					errors.Assert(t, test.err, err)
				}
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
				serv.validator.AssertCalled(t, "ValidateSchema", user)
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

func TestUpdate(t *testing.T) {
	mUser := mockUser()

	genReq := func(cb func(req *UpdateRequest)) *UpdateRequest {
		req := &UpdateRequest{}
		if cb != nil {
			cb(req)
		} else {
			req.Name = utils.NewString("New name")
		}
		return req
	}

	tests := []struct {
		name string
		id   string
		req  *UpdateRequest
		err  error
		mock func(s *mockService)
	}{{
		"invalid id",
		"abc123",
		genReq(nil),
		ErrNotFound,
		func(s *mockService) {
			s.repo.On("FindByID", "abc123").Return(nil, ErrRepositoryNotFound)
		},
	}, {
		"not enabled user",
		mUser.ID,
		genReq(nil),
		ErrNotFound,
		func(s *mockService) {
			u := copyUser(mUser)
			u.Enabled = false
			s.repo.On("FindByID", mUser.ID).Return(u, nil)
		},
	}, {
		"not validated user",
		mUser.ID,
		genReq(nil),
		ErrNotValidated,
		func(s *mockService) {
			u := copyUser(mUser)
			u.Validated = false
			s.repo.On("FindByID", mUser.ID).Return(u, nil)
		},
	}, {
		"invalid name and lastname",
		mUser.ID,
		genReq(func(req *UpdateRequest) {
			req.Name = utils.NewString("11111")
			req.Lastname = utils.NewString("22222")
		}),
		errors.Errors{ErrSchemaValidation},
		func(s *mockService) {
			u := copyUser(mUser)
			s.repo.On("FindByID", mUser.ID).Return(u, nil)
			s.validator.On("ValidateSchema", u).Return(ErrSchemaValidation)
		},
	}, {
		"invalid name and password",
		mUser.ID,
		genReq(func(req *UpdateRequest) {
			req.Name = utils.NewString("11111")
			req.Password = utils.NewString("222")
		}),
		errors.Errors{ErrPasswordValidation, ErrSchemaValidation},
		func(s *mockService) {
			u := copyUser(mUser)
			s.repo.On("FindByID", mUser.ID).Return(u, nil)
			s.validator.On("ValidatePassword", "222").Return(ErrPasswordValidation)
			s.validator.On("ValidateSchema", u).Return(ErrSchemaValidation)
		},
	}, {
		"invalid username and email",
		mUser.ID,
		genReq(func(req *UpdateRequest) {
			req.Username = utils.NewString("new user")
			req.Email = utils.NewString("inválid@-email.c")
		}),
		errors.Errors{ErrSchemaValidation},
		func(s *mockService) {
			u := copyUser(mUser)
			s.repo.On("FindByID", mUser.ID).Return(u, nil)
			s.validator.On("ValidateSchema", u).Return(ErrSchemaValidation)
		},
	}, {
		"invalid password",
		mUser.ID,
		genReq(func(req *UpdateRequest) {
			req.Password = utils.NewString("1245")
		}),
		errors.Errors{ErrPasswordValidation},
		func(s *mockService) {
			u := copyUser(mUser)
			s.repo.On("FindByID", mUser.ID).Return(u, nil)
			s.validator.On("ValidateSchema", mock.AnythingOfType("*models.User")).Return(nil)
			s.validator.On("ValidatePassword", "1245").Return(ErrPasswordValidation)
		},
	}, {
		"username not available",
		mUser.ID,
		genReq(func(req *UpdateRequest) {
			req.Username = utils.NewString("new-user")
		}),
		ErrNotAvailable,
		func(s *mockService) {
			u := copyUser(mUser)
			s.repo.On("FindByID", mUser.ID).Return(u, nil)
			eUser := mockUser()
			eUser.Username = "new-user"
			s.repo.On("FindByUsername", "new-user").Return(eUser, nil)

			s.validator.On("ValidateSchema", mock.AnythingOfType("*models.User")).Return(nil)
		},
	}, {
		"email not available",
		mUser.ID,
		genReq(func(req *UpdateRequest) {
			req.Email = utils.NewString("new@email.com")
		}),
		ErrNotAvailable,
		func(s *mockService) {
			u := copyUser(mUser)
			s.repo.On("FindByID", mUser.ID).Return(u, nil)
			eUser := mockUser()
			eUser.Email = "new@email.com"
			s.repo.On("FindByEmail", "new@email.com").Return(eUser, nil)

			s.validator.On("ValidateSchema", mock.AnythingOfType("*models.User")).Return(nil)
		},
	}, {
		"username and email not available",
		mUser.ID,
		genReq(func(req *UpdateRequest) {
			req.Username = utils.NewString("new-user")
			req.Email = utils.NewString("new@email.com")
		}),
		ErrNotAvailable.F("username", "not_available").F("email", "not_available"),
		func(s *mockService) {
			u := copyUser(mUser)
			s.repo.On("FindByID", mUser.ID).Return(u, nil)
			eu1 := mockUser()
			eu1.Username = "new-user"
			eu2 := mockUser()
			eu2.Email = "new@email.com"
			s.repo.On("FindByUsername", "new-user").Return(eu1, nil)
			s.repo.On("FindByEmail", "new@email.com").Return(eu2, nil)
			s.validator.On("ValidateSchema", mock.AnythingOfType("*models.User")).Return(nil)
		},
	}, {
		"error on update",
		mUser.ID,
		genReq(nil),
		ErrUpdate.Wrap(ErrRepositoryUpdate),
		func(s *mockService) {
			u := copyUser(mUser)
			s.repo.On("FindByID", mUser.ID).Return(u, nil)
			s.validator.On("ValidateSchema", mock.AnythingOfType("*models.User")).Return(nil)
			s.repo.On("Update", mock.AnythingOfType("*models.User")).Return(ErrRepositoryUpdate)
		},
	}, {
		"error on publishing",
		mUser.ID,
		genReq(nil),
		ErrUpdate.Wrap(events.ErrPublish),
		func(s *mockService) {
			u := copyUser(mUser)
			s.repo.On("FindByID", mUser.ID).Return(u, nil)
			s.validator.On("ValidateSchema", mock.AnythingOfType("*models.User")).Return(nil)
			s.repo.On("Update", mock.AnythingOfType("*models.User")).Return(nil)
			s.events.On("Publish", mock.AnythingOfType("*users.UserEvent"), mock.AnythingOfType("*events.Options")).Return(events.ErrPublish)
		},
	}, {
		"valid update",
		mUser.ID,
		genReq(func(req *UpdateRequest) {
			req.Username = utils.NewString("new-user")
			req.Email = utils.NewString("new@email.com")
			req.Password = utils.NewString("new-password")
		}),
		nil,
		func(s *mockService) {
			u := copyUser(mUser)
			s.repo.On("FindByID", mUser.ID).Return(u, nil)
			s.repo.On("FindByUsername", "new-user").Return(nil, ErrRepositoryNotFound)
			s.repo.On("FindByEmail", "new@email.com").Return(nil, ErrRepositoryNotFound)
			s.validator.On("ValidateSchema", mock.AnythingOfType("*models.User")).Return(nil)
			s.validator.On("ValidatePassword", "new-password").Return(nil)
			s.crypt.On("Hash", "new-password").Return("hashed.password", nil)
			s.repo.On("Update", mock.AnythingOfType("*models.User")).Return(nil)
			s.events.On("Publish", mock.AnythingOfType("*users.UserEvent"), mock.AnythingOfType("*events.Options")).Return(nil)
		},
	}, {
		"change name only",
		mUser.ID,
		genReq(nil),
		nil,
		func(s *mockService) {
			u := copyUser(mUser)
			s.repo.On("FindByID", mUser.ID).Return(u, nil)
			s.validator.On("ValidateSchema", mock.AnythingOfType("*models.User")).Return(nil)
			s.repo.On("Update", mock.AnythingOfType("*models.User")).Return(nil)
			s.events.On("Publish", mock.AnythingOfType("*users.UserEvent"), mock.AnythingOfType("*events.Options")).Return(nil)
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
				if assert.NotNil(err) {
					errors.Assert(t, test.err, err)
				}
				assert.Nil(user)
			} else {
				assert.Nil(err)
				if assert.NotNil(user) {
					assert.Equal(test.id, user.ID)
					assert.NotEqual(mUser, user)
					if test.req.Password != nil {
						assert.NotEqual(*test.req.Password, user.Password)
						assert.Equal("hashed.password", user.Password)
					}
				}
				serv.validator.AssertCalled(t, "ValidateSchema", user)
				serv.repo.AssertCalled(t, "Update", user)
			}
			serv.crypt.AssertExpectations(t)
			serv.repo.AssertExpectations(t)
			serv.validator.AssertExpectations(t)
			serv.events.AssertExpectations(t)
			serv.authServ.AssertExpectations(t)
		})
	}
}

func TestDelete(t *testing.T) {
	mUser := mockUser()

	tests := []struct {
		name string
		id   string
		err  error
		mock func(s *mockService)
	}{{
		"not found",
		"abc123",
		ErrNotFound,
		func(s *mockService) {
			s.repo.On("FindByID", "abc123").Return(nil, ErrRepositoryNotFound)
		},
	}, {
		"not enabled",
		mUser.ID,
		ErrNotFound,
		func(s *mockService) {
			u := copyUser(mUser)
			u.Enabled = false
			s.repo.On("FindByID", mUser.ID).Return(u, nil)
		},
	}, {
		"not validated",
		mUser.ID,
		ErrNotValidated,
		func(s *mockService) {
			u := copyUser(mUser)
			u.Validated = false
			s.repo.On("FindByID", mUser.ID).Return(u, nil)
		},
	}, {
		"error on delete",
		mUser.ID,
		ErrDelete.Wrap(ErrRepositoryDelete),
		func(s *mockService) {
			u := copyUser(mUser)
			s.repo.On("FindByID", mUser.ID).Return(u, nil)
			s.repo.On("Delete", mUser.ID).Return(ErrRepositoryDelete)
		},
	}, {
		"error on publish",
		mUser.ID,
		ErrDelete.Wrap(events.ErrPublish),
		func(s *mockService) {
			u := copyUser(mUser)
			s.repo.On("FindByID", mUser.ID).Return(u, nil)
			s.repo.On("Delete", mUser.ID).Return(nil)
			s.events.On("Publish", mock.AnythingOfType("*users.UserEvent"), mock.AnythingOfType("*events.Options")).Return(events.ErrPublish)
		},
	}, {
		"success",
		mUser.ID,
		nil,
		func(s *mockService) {
			u := copyUser(mUser)
			s.repo.On("FindByID", mUser.ID).Return(u, nil)
			s.repo.On("Delete", mUser.ID).Return(nil)
			s.events.On("Publish", mock.AnythingOfType("*users.UserEvent"), mock.AnythingOfType("*events.Options")).Return(nil)
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

func TestLogin(t *testing.T) {
	mUser := mockUser()
	mTokenStr := "encoded.token"

	genReq := func(cb func(req *LoginRequest)) *LoginRequest {
		req := &LoginRequest{
			UsernameOrEmail: utils.NewString("user"),
			Password:        utils.NewString("12345678"),
		}
		if cb != nil {
			cb(req)
		}
		return req
	}

	tests := []struct {
		name string
		req  *LoginRequest
		err  error
		mock func(s *mockService)
	}{{
		"empty request",
		genReq(func(req *LoginRequest) {
			req.UsernameOrEmail = nil
			req.Password = nil
		}),
		ErrInvalidLogin,
		nil,
	}, {
		"invalid username or email",
		genReq(func(req *LoginRequest) {
			req.UsernameOrEmail = nil
		}),
		ErrInvalidLogin,
		nil,
	}, {
		"invalid username",
		genReq(func(req *LoginRequest) {
			*req.UsernameOrEmail = "qwerty"
		}),
		ErrInvalidUser,
		func(s *mockService) {
			s.repo.On("FindByUsername", "qwerty").Return(nil, ErrRepositoryNotFound)
			s.repo.On("FindByEmail", "qwerty").Return(nil, ErrRepositoryNotFound)
		},
	}, {
		"invalid password",
		genReq(func(req *LoginRequest) {
			*req.Password = "wrong-password"
		}),
		ErrInvalidUser,
		func(s *mockService) {
			s.repo.On("FindByUsername", "user").Return(mUser, nil)
			s.crypt.On("Compare", mUser.Password, "wrong-password").Return(false)
		},
	}, {
		"login with username and password",
		genReq(nil),
		nil,
		func(s *mockService) {
			s.repo.On("FindByUsername", "user").Return(mUser, nil)
			s.crypt.On("Compare", mUser.Password, "12345678").Return(true)
			s.authServ.On("Create", mUser.ID).Return(mTokenStr, nil)
			// s.events.On("Publish", mock.Anything, mock.Anything).Return(nil)
		},
	}, {
		"login with email and password",
		genReq(func(req *LoginRequest) {
			*req.UsernameOrEmail = "user@user.com"
			*req.Password = "complexPassword#!"
		}),
		nil,
		func(s *mockService) {
			s.repo.On("FindByUsername", "user@user.com").Return(nil, ErrRepositoryNotFound)
			s.repo.On("FindByEmail", "user@user.com").Return(mUser, nil)
			s.crypt.On("Compare", mUser.Password, "complexPassword#!").Return(true)
			s.authServ.On("Create", mUser.ID).Return(mTokenStr, nil)
			// s.events.On("Publish", mock.Anything, mock.Anything).Return(nil)
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
				if assert.NotNil(err) {
					errors.Assert(t, test.err, err)
				}
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

func TestLogout(t *testing.T) {
	mUser := mockUser()
	mToken := models.NewToken(mUser.ID)
	mTokenStr := "encoded.token"

	tests := []struct {
		name     string
		tokenStr string
		err      error
		mock     func(s *mockService)
	}{{
		"empty token",
		"",
		ErrInvalidUser,
		func(s *mockService) {
			s.authServ.On("Invalidate", "").Return(nil, auth.ErrInvalidate)
		},
	}, {
		"non existing user",
		mTokenStr,
		ErrInvalidUser,
		func(s *mockService) {
			s.authServ.On("Invalidate", mTokenStr).Return(mToken, nil)
			s.repo.On("FindByID", mUser.ID).Return(nil, ErrRepositoryNotFound)
		},
	}, {
		"success",
		mTokenStr,
		nil,
		func(s *mockService) {
			s.authServ.On("Invalidate", mTokenStr).Return(mToken, nil)
			s.repo.On("FindByID", mUser.ID).Return(mUser, nil)
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			serv := newMockService()
			if test.mock != nil {
				test.mock(serv)
			}

			err := serv.Logout(test.tokenStr)

			if test.err != nil {
				if assert.NotNil(err) {
					errors.Assert(t, test.err, err)
				}
			} else {
				assert.Nil(err)
			}
			serv.repo.AssertExpectations(t)
			serv.authServ.AssertExpectations(t)
			serv.events.AssertExpectations(t)
			serv.validator.AssertExpectations(t)
			serv.crypt.AssertExpectations(t)
		})
	}
}
