package users

import (
	"reflect"
	"testing"

	"github.com/aboglioli/big-brother/errors"
	"github.com/aboglioli/big-brother/events"
	"github.com/aboglioli/big-brother/mock"
)

func TestGetByID(t *testing.T) {
	user1 := newMockUser()
	user2 := newMockUser()
	user2.Validated = true
	user3 := newMockUser()
	user3.Validated = true
	user3.Active = false

	tests := []struct {
		id   string
		err  error
		user *User
	}{{
		"123",
		ErrNotFound,
		user1,
	}, {
		"abc1235",
		ErrNotFound,
		user1,
	}, {
		"xyz123",
		ErrNotFound,
		user1,
	}, {
		user1.ID.Hex(),
		ErrNotValidated,
		user1,
	}, {
		user2.ID.Hex(),
		nil,
		user2,
	}, {
		user3.ID.Hex(),
		ErrNotActive,
		user3,
	}}

	for i, test := range tests {
		mockServ := newMockService()
		mockServ.repo.populate(test.user)
		u, err := mockServ.GetByID(test.id)

		if err != nil {
			if test.err != nil {
				expectedErr := test.err.(errors.Error)
				err := err.(errors.Error)
				if expectedErr.Code != err.Code {
					t.Errorf("test %d:\n-expected:%#v\n-actual:  %#v", i, test.err, err)
				}

			} else {
				t.Errorf("test %d:\n-expected: nil error\n-actual:%#v", i, err)
			}
		} else {
			if !reflect.DeepEqual(u, test.user) {
				t.Errorf("test %d:\n-expected:%#v\n-actual:  %#v", i, test.user, u)
			}
		}

		mockServ.repo.Mock.Assert(t,
			mock.Call("FindByID", test.id),
		)
	}
}

func TestRegister(t *testing.T) {
	user := newMockUser()
	user.Validated = true

	t.Run("Error", func(t *testing.T) {
		tests := []struct {
			req *RegisterRequest
			err error
		}{{
			&RegisterRequest{"admin", "1234567", "admin@admin.com"},
			errors.Errors{ErrPasswordValidation},
		}, {
			&RegisterRequest{"user", "123456789", "admin@admi.com"},
			errors.Errors{ErrNotAvailable},
		}, {
			&RegisterRequest{"admin", "12345678", "user@user.com"},
			errors.Errors{ErrNotAvailable},
		}, {
			&RegisterRequest{"user", "12345678", "user@user.com"},
			errors.Errors{ErrNotAvailable},
		}, {
			&RegisterRequest{"us€r", "1234", "user@user"},
			errors.Errors{ErrPasswordValidation, ErrSchemaValidation},
		}, {
			&RegisterRequest{"user", "1234", "user@user"},
			errors.Errors{ErrNotAvailable, ErrPasswordValidation, ErrSchemaValidation},
		}}

		for i, test := range tests {
			mockServ := newMockService()
			existingUser := newMockUser()
			mockServ.repo.populate(existingUser)
			u, err := mockServ.Register(test.req)

			if err == nil || u != nil {
				t.Errorf("test %d: expected error and nil user", i)
			}

			if !errors.Compare(test.err, err) {
				t.Errorf("test %d:\n-expected:%#v\n-actual:  %#v", i, test.err, err)
			}

			call1 := mock.Call("FindByUsername", test.req.Username)
			call2 := mock.Call("FindByEmail", test.req.Email)

			if test.req.Username == existingUser.Username {
				call1 = call1.Return(mock.NotNil, mock.Nil)
			} else {
				call1 = call1.Return(mock.Nil, mock.NotNil)
			}

			if test.req.Email == existingUser.Email {
				call2 = call2.Return(mock.NotNil, mock.Nil)
			} else {
				call2 = call2.Return(mock.Nil, mock.NotNil)
			}

			mockServ.repo.Mock.Assert(t,
				call1,
				call2,
			)

			mockServ.validator.Mock.Assert(t,
				mock.Call("ValidatePassword", test.req.Password),
				mock.Call("ValidateSchema", mock.NotNil),
			)
		}
	})

	t.Run("OK", func(t *testing.T) {
		tests := []struct {
			req  *RegisterRequest
			user *User
		}{{
			&RegisterRequest{"admin", "123456789", "admin@admin.com"},
			&User{
				Username: "admin",
				Email:    "admin@admin.com",
			},
		}, {
			&RegisterRequest{"user", "asdqwe123", "user@user.com"},
			&User{
				Username: "user",
				Email:    "user@user.com",
			},
		}}

		for i, test := range tests {
			mockServ := newMockService()
			user, err := mockServ.Register(test.req)

			// Response
			if err != nil || user == nil {
				t.Errorf("test %d: expected user, got error", i)
				continue
			}

			// Properties
			if user.Username != test.user.Username || user.Email != test.user.Email {
				t.Errorf("test %d:\n-expected:%s - %s\n-actual:  %s - %s", i, test.user.Username, test.user.Email, user.Username, user.Email)
			}

			if user.Password == test.req.Password || len(user.Password) < 10 {
				t.Errorf("test %d: password wrong hashing: %s", i, user.Password)
			}

			if !user.Enabled || !user.Active || user.Validated {
				t.Errorf("test %d: %v - %v - %v", i, user.Enabled, user.Active, user.Validated)
			}

			// Validator
			mockServ.validator.Mock.Assert(t,
				mock.Call("ValidatePassword", test.req.Password).Return(nil),
				mock.Call("ValidateSchema", user).Return(nil),
			)

			// Repository
			mockServ.repo.Mock.Assert(t,
				mock.Call("FindByUsername", test.req.Username).Return(mock.Nil, mock.NotNil),
				mock.Call("FindByEmail", test.req.Email).Return(mock.Nil, mock.NotNil),
				mock.Call("Insert", mock.NotNil).Return(nil),
			)
			insertedUser, ok := mockServ.repo.Mock.Calls[2].Args[0].(*User)
			if !ok {
				t.Error("invalid conversion")
				continue
			}
			if !reflect.DeepEqual(user, insertedUser) {
				t.Errorf("test %d: inserted user is not equal to returned user\n-expected:%#v\n-actual:  %#v", i, user, insertedUser)
			}
			insertedUser = mockServ.repo.Collection[0]
			if !reflect.DeepEqual(user, insertedUser) {
				t.Errorf("test %d: inserted user is not equal to returned user\n-expected:%#v\n-actual:  %#v", i, user, insertedUser)
			}

			// Events
			mockServ.events.Mock.Assert(t,
				mock.Call("Publish", mock.NotNil, mock.NotNil).Return(nil),
			)
			event, ok1 := mockServ.events.Mock.Calls[0].Args[0].(*UserEvent)
			opts, ok2 := mockServ.events.Mock.Calls[0].Args[1].(*events.Options)

			if !ok1 || !ok2 {
				t.Error("invalid conversion")
				continue
			}
			if event.Type != "UserCreated" {
				t.Errorf("test %d: invalid event type", i)
			}
			if !reflect.DeepEqual(user, event.User) {
				t.Errorf("test %d:\n-expected:%#v\n-actual:  %#v", i, user, event.User)
			}
			if opts.Exchange != "user" || opts.Route != "user.created" || opts.Queue != "" {
				t.Errorf("test %d: invalid event options %#v", i, opts)
			}
		}
	})
}

func TestUpdate(t *testing.T) {
	user1 := newMockUser()
	user1.Validated = true
	user2 := newMockUser()
	user2.Validated = true
	user2.Username = "admin"
	user2.Email = "admin@admin.com"
	user2.Name = "Admin"
	user2.Lastname = "Admin"
	user2.Roles = []Role{ADMIN}

	t.Run("Error", func(t *testing.T) {
		req1 := &UpdateRequest{
			Name:     new(string),
			Lastname: new(string),
		}
		*req1.Name = "11111"
		*req1.Lastname = "22222"
		req2 := &UpdateRequest{
			Username: new(string),
		}
		*req2.Username = "admin"
		req3 := &UpdateRequest{
			Username: new(string),
		}
		*req3.Username = "user#1"
		req4 := &UpdateRequest{
			Password: new(string),
		}
		*req4.Password = "1234"
		req5 := &UpdateRequest{
			Email: new(string),
		}
		*req5.Email = "admin@admin"
		req6 := &UpdateRequest{
			Email: new(string),
		}
		*req6.Email = "admin@admin.com"
		req7 := &UpdateRequest{
			Username: new(string),
			Email:    new(string),
			Password: new(string),
			Name:     new(string),
		}
		*req7.Username = "admin"
		*req7.Email = "admin@admin.com"
		*req7.Password = "1234"
		*req7.Name = "Al@n"

		tests := []struct {
			id  string
			req *UpdateRequest
			err error
		}{{
			"123",
			&UpdateRequest{},
			ErrNotFound,
		}, {
			"123456",
			req1,
			ErrNotFound,
		}, {
			user1.ID.Hex(),
			req1,
			errors.Errors{ErrSchemaValidation},
		}, {
			user1.ID.Hex(),
			req2,
			errors.Errors{ErrNotAvailable},
		}, {
			user1.ID.Hex(),
			req3,
			errors.Errors{ErrSchemaValidation},
		}, {
			user1.ID.Hex(),
			req4,
			errors.Errors{ErrPasswordValidation},
		}, {
			user1.ID.Hex(),
			req5,
			errors.Errors{ErrSchemaValidation},
		}, {
			user1.ID.Hex(),
			req6,
			errors.Errors{ErrNotAvailable},
		}, {
			user1.ID.Hex(),
			req7,
			errors.Errors{ErrNotAvailable, ErrPasswordValidation, ErrSchemaValidation},
		}}

		for i, test := range tests {
			mockServ := newMockService()
			mockServ.repo.populate(user1, user2)
			user, err := mockServ.Update(test.id, test.req)

			if user != nil || err == nil {
				t.Errorf("test %d: expected nil user and error\ngot: %#v - %#v", i, user, err)
				continue
			}

			if !errors.Compare(err, test.err) {
				t.Errorf("test %d:\n-expected:%#v\n-actual:  %#v", i, test.err, err)
			}

			repoCalls := mock.Calls{mock.Call("FindByID", test.id)}
			validatorCalls := mock.Calls{}

			if test.id == user1.ID.Hex() {
				if test.req.Username != nil {
					repoCalls = append(repoCalls, mock.Call("FindByUsername", *test.req.Username))
				}
				if test.req.Password != nil {
					validatorCalls = append(validatorCalls, mock.Call("ValidatePassword", *test.req.Password))
				}
				if test.req.Email != nil {
					repoCalls = append(repoCalls, mock.Call("FindByEmail", *test.req.Email))
				}

				validatorCalls = append(validatorCalls, mock.Call("ValidateSchema", mock.NotNil))
			}

			mockServ.validator.Mock.Assert(t,
				validatorCalls...,
			)

			mockServ.repo.Mock.Assert(t,
				repoCalls...,
			)
		}
	})

	t.Run("OK", func(t *testing.T) {
		req1 := &UpdateRequest{
			Name:     new(string),
			Lastname: new(string),
		}
		*req1.Name = "NewName"
		*req1.Lastname = "NewLastname"

		tests := []struct {
			id   string
			req  *UpdateRequest
			user *User
		}{{
			user1.ID.Hex(),
			req1,
			&User{
				Username: "user",
				Password: "123456789",
				Email:    "user@user.com",
				Name:     "NewName",
				Lastname: "NewLastname",
			},
		}}

		for i, test := range tests {
			mockServ := newMockService()
			mockServ.repo.populate(user1, user2)
			user, err := mockServ.Update(test.id, test.req)

			if user == nil || err != nil {
				t.Errorf("test %d: expected user with id %s, got error: %#v", i, user.ID.Hex(), err)
				continue
			}

			if user.Username != test.user.Username || user.Email != test.user.Email || user.Name != test.user.Name || user.Lastname != test.user.Lastname {
				t.Errorf("test %d:\n-expected:%#v\n-actual:  %#v", i, test.user, user)
			}

			if !user.ComparePassword(test.user.Password) {
				t.Errorf("test %d: password does not match", i)
			}

			repoCalls := mock.Calls{mock.Call("FindByID", test.id)}
			validatorCalls := mock.Calls{}
			if test.req.Username != nil {
				repoCalls = append(repoCalls, mock.Call("FindByUsername", *test.req.Username))
			}
			if test.req.Password != nil {
				validatorCalls = append(validatorCalls, mock.Call("ValidatePassword", test.req.Password))
			}
			if test.req.Email != nil {
				repoCalls = append(repoCalls, mock.Call("FindByEmail", *test.req.Email))
			}
			repoCalls = append(repoCalls, mock.Call("Update", user))
			validatorCalls = append(validatorCalls, mock.Call("ValidateSchema", user))

			mockServ.repo.Mock.Assert(t, repoCalls...)
			mockServ.validator.Mock.Assert(t, validatorCalls...)
			mockServ.events.Mock.Assert(t, mock.Call("Publish", mock.NotNil, mock.NotNil))

			// Events
			mockServ.events.Mock.Assert(t,
				mock.Call("Publish", mock.NotNil, mock.NotNil).Return(nil),
			)
			event, ok1 := mockServ.events.Mock.Calls[0].Args[0].(*UserEvent)
			opts, ok2 := mockServ.events.Mock.Calls[0].Args[1].(*events.Options)

			if !ok1 || !ok2 {
				t.Error("invalid conversion")
				continue
			}
			if event.Type != "UserUpdated" {
				t.Errorf("test %d: invalid event type", i)
			}
			if !reflect.DeepEqual(user, event.User) {
				t.Errorf("test %d:\n-expected:%#v\n-actual:  %#v", i, user, event.User)
			}
			if opts.Exchange != "user" || opts.Route != "user.updated" || opts.Queue != "" {
				t.Errorf("test %d: invalid event options %#v", i, opts)
			}
		}
	})
}

func TestDelete(t *testing.T) {
	user1 := newMockUser()
	user1.Validated = true
	user2 := newMockUser()
	user2.Validated = true
	user2.Username = "admin"
	user2.Email = "admin@admin.com"
	user2.Name = "Admin"
	user2.Lastname = "Admin"
	user2.Roles = []Role{ADMIN}

	t.Run("Error", func(t *testing.T) {
		tests := []struct {
			id  string
			err error
		}{{
			"1234",
			ErrNotFound,
		}}

		for i, test := range tests {
			mockServ := newMockService()
			mockServ.repo.populate(user1, user2)
			err := mockServ.Delete(test.id)

			if !errors.Compare(err, test.err) {
				t.Errorf("test %d:\n-expected:%#v\n-actual:  %#v", i, test.err, err)
			}
		}
	})

	t.Run("OK", func(t *testing.T) {
		tests := []struct {
			id string
		}{{
			user1.ID.Hex(),
		}, {
			user2.ID.Hex(),
		}}

		for i, test := range tests {
			mockServ := newMockService()
			mockServ.repo.populate(user1, user2)
			err := mockServ.Delete(test.id)

			if err != nil {
				t.Errorf("test %d: error not expected %#v", i, err)
			}
		}
	})

}
