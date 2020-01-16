package users

import (
	"testing"

	"github.com/aboglioli/big-brother/internal/auth"
	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/events"
	"github.com/aboglioli/big-brother/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func mockUser() *User {
	user := NewUser()
	user.Username = "user"
	user.SetPassword("12345678")
	user.Email = "user@user.com"
	user.Name = "Name"
	user.Lastname = "Lastname"
	user.Validated = true
	user.Enabled = true
	return user
}

func copyUser(u *User) *User {
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
		mUser.ID.Hex(),
		ErrNotFound,
		func(s *mockService) {
			s.repo.On("FindByID", mUser.ID.Hex()).Return(nil, ErrRepositoryNotFound)
		},
	}, {
		"not validated",
		mUser.ID.Hex(),
		ErrNotValidated,
		func(s *mockService) {
			u := copyUser(mUser)
			u.Validated = false
			s.repo.On("FindByID", mUser.ID.Hex()).Return(u, nil)
		},
	}, {
		"not enabled",
		mUser.ID.Hex(),
		ErrNotFound,
		func(s *mockService) {
			u := copyUser(mUser)
			u.Enabled = false
			s.repo.On("FindByID", mUser.ID.Hex()).Return(u, nil)
		},
	}, {
		"existing user",
		mUser.ID.Hex(),
		nil,
		func(s *mockService) {
			s.repo.On("FindByID", mUser.ID.Hex()).Return(mUser, nil)
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
					assert.Equal(test.id, user.ID.Hex())
				}
			}
			serv.repo.AssertExpectations(t)
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
		errors.Errors{ErrSchemaValidation, ErrPasswordValidation},
		func(s *mockService) {
			s.validator.On("ValidateSchema", mock.Anything).Return(ErrSchemaValidation)
			s.validator.On("ValidatePassword", "12345678").Return(ErrPasswordValidation)
		},
	}, {
		"invalid password",
		genReq(func(req *RegisterRequest) {
			req.Password = "1234567"
		}),
		errors.Errors{ErrPasswordValidation},
		func(s *mockService) {
			s.validator.On("ValidateSchema", mock.Anything).Return(nil)
			s.validator.On("ValidatePassword", "1234567").Return(ErrPasswordValidation)
		},
	}, {
		"username not available",
		genReq(nil),
		errors.Errors{ErrNotAvailable},
		func(s *mockService) {
			s.repo.On("FindByUsername", "user").Return(mUser, nil)
			s.repo.On("FindByEmail", "user@user.com").Return(nil, ErrRepositoryNotFound)
			s.validator.On("ValidateSchema", mock.Anything).Return(nil)
			s.validator.On("ValidatePassword", "12345678").Return(nil)
		},
	}, {
		"email not available",
		genReq(nil),
		errors.Errors{ErrNotAvailable},
		func(s *mockService) {
			s.repo.On("FindByUsername", "user").Return(nil, ErrRepositoryNotFound)
			s.repo.On("FindByEmail", "user@user.com").Return(mUser, nil)
			s.validator.On("ValidateSchema", mock.Anything).Return(nil)
			s.validator.On("ValidatePassword", "12345678").Return(nil)
		},
	}, {
		"user and email not available",
		genReq(nil),
		errors.Errors{ErrNotAvailable},
		func(s *mockService) {
			s.repo.On("FindByUsername", "user").Return(mUser, nil)
			s.repo.On("FindByEmail", "user@user.com").Return(mUser, nil)
			s.validator.On("ValidateSchema", mock.Anything).Return(nil)
			s.validator.On("ValidatePassword", "12345678").Return(nil)
		},
	}, {
		"error on insert",
		genReq(nil),
		ErrRegister,
		func(s *mockService) {
			s.repo.On("FindByUsername", "user").Return(nil, ErrRepositoryNotFound)
			s.repo.On("FindByEmail", "user@user.com").Return(nil, ErrRepositoryNotFound)
			s.validator.On("ValidateSchema", mock.Anything).Return(nil)
			s.validator.On("ValidatePassword", "12345678").Return(nil)
			s.repo.On("Insert", mock.Anything).Return(ErrRepositoryInsert)
		},
	}, {
		"error on publishing event",
		genReq(nil),
		ErrRegister,
		func(s *mockService) {
			s.repo.On("FindByUsername", "user").Return(nil, ErrRepositoryNotFound)
			s.repo.On("FindByEmail", "user@user.com").Return(nil, ErrRepositoryNotFound)
			s.validator.On("ValidateSchema", mock.Anything).Return(nil)
			s.validator.On("ValidatePassword", "12345678").Return(nil)
			s.repo.On("Insert", mock.Anything).Return(nil)
			s.events.On("Publish", mock.Anything, mock.Anything).Return(events.ErrPublish)
		},
	}, {
		"valid user",
		genReq(nil),
		nil,
		func(s *mockService) {
			s.repo.On("FindByUsername", "user").Return(nil, ErrRepositoryNotFound)
			s.repo.On("FindByEmail", "user@user.com").Return(nil, ErrRepositoryNotFound)
			s.validator.On("ValidateSchema", mock.Anything).Return(nil)
			s.validator.On("ValidatePassword", "12345678").Return(nil)
			s.repo.On("Insert", mock.Anything).Return(nil)
			s.events.On("Publish", mock.Anything, mock.Anything).Return(nil)
		},
	}, {
		"valid admin",
		genReq(func(req *RegisterRequest) {
			req.Username = "admin"
			req.Password = "123456789"
			req.Email = "admin@admin.com"
			req.Name = "Admin"
			req.Lastname = "Lastname"
		}),
		nil,
		func(s *mockService) {
			s.repo.On("FindByUsername", "admin").Return(nil, ErrRepositoryNotFound)
			s.repo.On("FindByEmail", "admin@admin.com").Return(nil, ErrRepositoryNotFound)
			s.validator.On("ValidateSchema", mock.Anything).Return(nil)
			s.validator.On("ValidatePassword", "123456789").Return(nil)
			s.repo.On("Insert", mock.Anything).Return(nil)
			s.events.On("Publish", mock.Anything, mock.Anything).Return(nil)
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
					assert.NotEmpty(user.ID.Hex())
					assert.Equal(test.req.Username, user.Username)
					assert.NotEqual(test.req.Password, user.Password)
					assert.True(user.ComparePassword(test.req.Password))
					assert.Equal(test.req.Email, user.Email)
					assert.Equal(test.req.Name, user.Name)
					assert.Equal(test.req.Lastname, user.Lastname)
				}
				serv.validator.AssertCalled(t, "ValidateSchema", user)
				serv.repo.AssertCalled(t, "Insert", user)
			}
			serv.repo.AssertExpectations(t)
			serv.validator.AssertExpectations(t)
			serv.events.AssertExpectations(t)
		})
	}
}

func TestUpdate(t *testing.T) {
	mUser := mockUser()

	genReq := func(cb func(req *UpdateRequest)) *UpdateRequest {
		req := &UpdateRequest{
			Name: utils.NewString("New name"),
		}
		if cb != nil {
			cb(req)
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
		mUser.ID.Hex(),
		genReq(nil),
		ErrNotFound,
		func(s *mockService) {
			u := copyUser(mUser)
			u.Enabled = false
			s.repo.On("FindByID", mUser.ID.Hex()).Return(u, nil)
		},
	}, {
		"not validated user",
		mUser.ID.Hex(),
		genReq(nil),
		ErrNotValidated,
		func(s *mockService) {
			u := copyUser(mUser)
			u.Validated = false
			s.repo.On("FindByID", mUser.ID.Hex()).Return(u, nil)
		},
	}, {
		"invalid name and lastname",
		mUser.ID.Hex(),
		genReq(func(req *UpdateRequest) {
			*req.Name = "11111"
			*req.Lastname = "22222"
		}),
		ErrSchemaValidation,
		func(s *mockService) {
			u := copyUser(mUser)
			s.repo.On("FindByID", mUser.ID.Hex()).Return(u, nil)
			s.validator.On("ValidateSchema", u).Return(ErrSchemaValidation)
		},
	}, {
		"invalid username and email",
		mUser.ID.Hex(),
		genReq(func(req *UpdateRequest) {
			*req.Username = "new user"
			*req.Email = "inválid@-email.c"
		}),
		ErrSchemaValidation,
		func(s *mockService) {
			u := copyUser(mUser)
			s.repo.On("FindByID", mUser.ID.Hex()).Return(u, nil)
			s.validator.On("ValidateSchema", u).Return(ErrSchemaValidation)
		},
	}, {
		"invalid password",
		mUser.ID.Hex(),
		genReq(func(req *UpdateRequest) {
			*req.Password = "1245"
		}),
		ErrPasswordValidation,
		func(s *mockService) {
			u := copyUser(mUser)
			s.repo.On("FindByID", mUser.ID.Hex()).Return(u, nil)
			s.validator.On("ValidatePassword", "1245").Return(ErrPasswordValidation)
		},
	}, {
		"username not available",
		mUser.ID.Hex(),
		genReq(func(req *UpdateRequest) {
			*req.Username = "new-user"
		}),
		errors.Errors{ErrNotAvailable},
		func(s *mockService) {
			u := copyUser(mUser)
			s.repo.On("FindByID", mUser.ID.Hex()).Return(u, nil)
			eUser := mockUser()
			eUser.Username = "new-user"
			s.repo.On("FindByUsername", "new-user").Return(eUser, nil)
		},
	}, {
		"email not available",
		mUser.ID.Hex(),
		genReq(func(req *UpdateRequest) {
			*req.Email = "new@email.com"
		}),
		errors.Errors{ErrNotAvailable},
		func(s *mockService) {
			u := copyUser(mUser)
			s.repo.On("FindByID", mUser.ID.Hex()).Return(u, nil)
			eUser := mockUser()
			eUser.Email = "new@email.com"
			s.repo.On("FindByEmail", "new@email.com").Return(eUser, nil)
		},
	}, {
		"username and email not available",
		mUser.ID.Hex(),
		genReq(func(req *UpdateRequest) {
			*req.Username = "new-user"
			*req.Email = "new@email.com"
		}),
		errors.Errors{ErrNotAvailable},
		func(s *mockService) {
			u := copyUser(mUser)
			s.repo.On("FindByID", mUser.ID.Hex()).Return(u, nil)
			eu1 := mockUser()
			eu1.Username = "new-user"
			eu2 := mockUser()
			eu2.Email = "new@email.com"
			s.repo.On("FindByUsername", "new-user").Return(eu1, nil)
			s.repo.On("FindByEmail", "new@email.com").Return(eu2, nil)
		},
	}, {
		"valid update",
		mUser.ID.Hex(),
		genReq(func(req *UpdateRequest) {
			*req.Username = "new-user"
			*req.Email = "new@email.com"
			*req.Password = "new-password"
		}),
		nil,
		func(s *mockService) {
			u := copyUser(mUser)
			s.repo.On("FindByID", mUser.ID.Hex()).Return(u, nil)
			s.repo.On("FindByUsername", "new-user").Return(nil, ErrRepositoryNotFound)
			s.repo.On("FindByEmail", "new@email.com").Return(nil, ErrRepositoryNotFound)
			s.validator.On("ValidateSchema", mock.Anything).Return(nil)
			s.validator.On("ValidatePassword", "new-password").Return(nil)
			s.repo.On("Update", mock.Anything).Return(nil)
			s.events.On("Publish", mock.Anything, mock.Anything).Return(nil)
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
					assert.Equal(test.id, user.ID.Hex())
					assert.NotEqual(mUser, user)
				}
				serv.validator.AssertCalled(t, "ValidateSchema", user)
				serv.repo.AssertCalled(t, "Update", user)
			}
			serv.validator.AssertExpectations(t)
			serv.repo.AssertExpectations(t)
			serv.events.AssertExpectations(t)
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
		mUser.ID.Hex(),
		ErrNotFound,
		func(s *mockService) {
			u := copyUser(mUser)
			u.Enabled = false
			s.repo.On("FindByID", mUser.ID.Hex()).Return(u, nil)
		},
	}, {
		"not validated",
		mUser.ID.Hex(),
		ErrNotValidated,
		func(s *mockService) {
			u := copyUser(mUser)
			u.Validated = false
			s.repo.On("FindByID", mUser.ID.Hex()).Return(u, nil)
		},
	}, {
		"success",
		mUser.ID.Hex(),
		nil,
		func(s *mockService) {
			u := copyUser(mUser)
			s.repo.On("FindByID", mUser.ID.Hex()).Return(u, nil)
			s.repo.On("Delete", mUser.ID.Hex()).Return(nil)
			s.events.On("Publish", mock.Anything, mock.Anything).Return(nil)
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
			}
		})
	}
}

func TestLogin(t *testing.T) {
	mUser := mockUser()
	mToken := auth.NewToken(mUser.ID.Hex())

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
		},
	}, {
		"invalid password",
		genReq(func(req *LoginRequest) {
			*req.Password = "wrong-password"
		}),
		ErrInvalidUser,
		func(s *mockService) {
			s.repo.On("FindByUsername", "user").Return(mUser, nil)
		},
	}, {
		"login with username and password",
		genReq(nil),
		nil,
		func(s *mockService) {
			s.repo.On("FindByUsername", "user").Return(mUser, nil)
			s.authServ.On("Create", mUser.ID.Hex()).Return(mToken, nil)
			s.events.On("Publish", mock.Anything, mock.Anything).Return(nil)
		},
	}, {
		"login with email and password",
		genReq(func(req *LoginRequest) {
			*req.UsernameOrEmail = "user@user.com"
		}),
		nil,
		func(s *mockService) {
			s.repo.On("FindByUsername", "user@user.com").Return(nil, ErrRepositoryNotFound)
			s.repo.On("FindByEmail", "user@user.com").Return(mUser, nil)
			s.authServ.On("Create", mUser.ID.Hex()).Return(mToken, nil)
			s.events.On("Publish", mock.Anything, mock.Anything).Return(nil)
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			serv := newMockService()
			if test.mock != nil {
				test.mock(serv)
			}

			token, err := serv.Login(test.req)

			if test.err != nil {
				if assert.NotNil(err) {
					errors.Assert(t, test.err, err)
				}
				assert.Nil(token)
			} else {
				assert.Nil(err)
				if assert.NotNil(token) {
					assert.Equal(mToken, token)
				}
			}
			serv.repo.AssertExpectations(t)
			serv.authServ.AssertExpectations(t)
			serv.events.AssertExpectations(t)
		})
	}
}

func TestLogout(t *testing.T) {
	mUser := mockUser()
	mToken := auth.NewToken(mUser.ID.Hex())
	mTokenStr, err := mToken.Encode()
	require.Nil(t, err)

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
			s.repo.On("FindByID", mUser.ID.Hex()).Return(nil, ErrRepositoryNotFound)
		},
	}, {
		"success",
		mTokenStr,
		nil,
		func(s *mockService) {
			s.authServ.On("Invalidate", mTokenStr).Return(mToken, nil)
			s.repo.On("FindByID", mUser.ID.Hex()).Return(mUser, nil)
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
			serv.authServ.AssertExpectations(t)
			serv.repo.AssertExpectations(t)
		})
	}
}
