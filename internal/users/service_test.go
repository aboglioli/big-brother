package users

import (
	"reflect"
	"testing"

	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/events"
	"github.com/aboglioli/big-brother/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func TestGetByID(t *testing.T) {
	assert := assert.New(t)

	user1 := newMockUser("")
	user1.Validated = false
	user2 := newMockUser("")
	user2.Validated = true
	user3 := newMockUser("")
	user3.Enabled = false

	tests := []struct {
		id   string
		err  error
		user *User
	}{{
		"123",
		ErrNotFound,
		nil,
	}, {
		"abc1235",
		ErrNotFound,
		nil,
	}, {
		"xyz123",
		ErrNotFound,
		nil,
	}, {
		user1.ID.Hex(),
		ErrNotValidated,
		nil,
	}, {
		user2.ID.Hex(),
		nil,
		user2,
	}, {
		user3.ID.Hex(),
		ErrNotFound,
		nil,
	}}

	for i, test := range tests {
		mockServ := newMockService()
		mockServ.repo.populate(user1, user2)
		u, err := mockServ.GetByID(test.id)

		if !assert.True(errors.Compare(test.err, err), i) {
			t.Errorf("test %d:\n-expected:%s\n-actual:  %s", i, test.err, err)
		}
		assert.Equal(test.user, u, "test %d", i)

		if msg := mockServ.repo.Mock.Assert(
			mock.Call("FindByID", test.id),
		); msg != "" {
			t.Errorf("test %d: %s", i, msg)
		}
	}
}

func TestRegister(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		assert := assert.New(t)

		tests := []struct {
			req *RegisterRequest
			err error
		}{{
			&RegisterRequest{"admin", "1234567", "admin@admin.com", "Name", "Lastname"},
			errors.Errors{ErrPasswordValidation},
		}, {
			&RegisterRequest{"user", "123456789", "admin@admi.com", "Name", "Lastname"},
			errors.Errors{ErrNotAvailable},
		}, {
			&RegisterRequest{"admin", "12345678", "user@user.com", "Name", "Lastname"},
			errors.Errors{ErrNotAvailable},
		}, {
			&RegisterRequest{"user", "12345678", "user@user.com", "Name", "Lastname"},
			errors.Errors{ErrNotAvailable},
		}, {
			&RegisterRequest{"us€r", "1234", "user@user", "Name", "Lastname"},
			errors.Errors{ErrPasswordValidation, ErrSchemaValidation},
		}, {
			&RegisterRequest{"user", "1234", "user@user", "Name", "Lastname"},
			errors.Errors{ErrNotAvailable, ErrPasswordValidation, ErrSchemaValidation},
		}, {
			&RegisterRequest{"new-user", "123456789", "user@new-user.com", "Alan1", "Lastname2"},
			errors.Errors{ErrSchemaValidation},
		}, {
			&RegisterRequest{"user", "123456789", "user@new-user.com", "Alan1", "Lastname"},
			errors.Errors{ErrNotAvailable, ErrSchemaValidation},
		}}

		for i, test := range tests {
			mockServ := newMockService()
			existingUser := newMockUser("")
			mockServ.repo.populate(existingUser)
			u, err := mockServ.Register(test.req)

			if assert.NotNil(err, i) {
				assert.True(errors.Compare(test.err, err), i)
			}
			assert.Nil(u, i)

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

			if msg := mockServ.repo.Mock.Assert(
				call1,
				call2,
			); msg != "" {
				t.Errorf("test %d: %s", i, msg)
			}

			if msg := mockServ.validator.Mock.Assert(
				mock.Call("ValidatePassword", test.req.Password),
				mock.Call("ValidateSchema", mock.NotNil),
			); msg != "" {
				t.Errorf("test %d: %s", i, msg)
			}
		}
	})

	t.Run("OK", func(t *testing.T) {
		assert := assert.New(t)

		tests := []struct {
			req  *RegisterRequest
			user *User
		}{{
			&RegisterRequest{"admin", "123456789", "admin@admin.com", "Admin", "Lastname"},
			&User{
				Username: "admin",
				Email:    "admin@admin.com",
			},
		}, {
			&RegisterRequest{"user", "asdqwe123", "user@user.com", "User", "Lastname"},
			&User{
				Username: "user",
				Email:    "user@user.com",
			},
		}}

		for i, test := range tests {
			mockServ := newMockService()
			user, err := mockServ.Register(test.req)

			// Response
			assert.Nil(err, i)
			if assert.NotNil(user, i) {
				assert.Equal(test.user.Username, user.Username, i)
				assert.Equal(test.user.Email, user.Email)
				assert.NotEqual(test.req.Password, user.Password, i)
				assert.True(user.Enabled, i)
				assert.False(user.Validated, i)
			}

			// Validator
			if msg := mockServ.validator.Mock.Assert(
				mock.Call("ValidatePassword", test.req.Password).Return(nil),
				mock.Call("ValidateSchema", user).Return(nil),
			); msg != "" {
				t.Errorf("test %d: %s", i, msg)
			}

			// Repository
			if msg := mockServ.repo.Mock.Assert(
				mock.Call("FindByUsername", test.req.Username).Return(mock.Nil, mock.NotNil),
				mock.Call("FindByEmail", test.req.Email).Return(mock.Nil, mock.NotNil),
				mock.Call("Insert", mock.NotNil).Return(nil),
			); msg != "" {
				t.Errorf("test %d: %s", i, msg)
			}

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
			if msg := mockServ.events.Mock.Assert(
				mock.Call("Publish", mock.NotNil, mock.NotNil).Return(nil),
			); msg != "" {
				t.Errorf("test %d: %s", i, msg)
			}
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
	user1 := newMockUser("")
	user1.Validated = true
	user2 := newMockUser("")
	user2.Validated = true
	user2.Username = "admin"
	user2.Email = "admin@admin.com"
	user2.Name = "Admin"
	user2.Lastname = "Admin"
	user2.Roles = []Role{ADMIN}
	user3 := newMockUser("")
	user3.Validated = false
	user3.Username = "non-validated"
	user3.Email = "non-validated@user.com"
	user4 := newMockUser("")
	user4.Validated = true
	user4.Enabled = false
	user4.Username = "non-enabled"
	user4.Email = "non-enabled@user.com"

	t.Run("Error", func(t *testing.T) {
		assert := assert.New(t)

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
		req8 := &UpdateRequest{
			Name:     new(string),
			Lastname: new(string),
		}
		*req8.Name = "Hello"
		*req8.Lastname = "World"

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
		}, {
			user3.ID.Hex(),
			req8,
			ErrNotValidated,
		}, {
			user4.ID.Hex(),
			req8,
			ErrNotFound,
		}}

		for i, test := range tests {
			mockServ := newMockService()
			mockServ.repo.populate(user1, user2, user3)
			user, err := mockServ.Update(test.id, test.req)

			if assert.NotNil(err, i) {
				assert.True(errors.Compare(test.err, err), i)
			}
			assert.Nil(user, i)

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

			if msg := mockServ.validator.Mock.Assert(validatorCalls...); msg != "" {
				t.Errorf("test %d: %s", i, msg)
			}

			if msg := mockServ.repo.Mock.Assert(repoCalls...); msg != "" {
				t.Errorf("test %d: %s", i, msg)
			}
		}
	})

	t.Run("OK", func(t *testing.T) {
		assert := assert.New(t)

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

			assert.Nil(err, i)
			if assert.NotNil(user, i) {
				assert.Equal(test.user.Username, user.Username, i)
				assert.Equal(test.user.Email, user.Email, i)
				assert.Equal(test.user.Name, user.Name, i)
				assert.Equal(test.user.Lastname, user.Lastname, i)
				assert.True(user.ComparePassword(test.user.Password), i)
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

			if msg := mockServ.repo.Mock.Assert(repoCalls...); msg != "" {
				t.Errorf("test %d: %s", i, msg)
			}
			if msg := mockServ.validator.Mock.Assert(validatorCalls...); msg != "" {
				t.Errorf("test %d: %s", i, msg)
			}
			if msg := mockServ.events.Mock.Assert(mock.Call("Publish", mock.NotNil, mock.NotNil)); msg != "" {
				t.Errorf("test %d: %s", i, msg)
			}

			// Events
			if msg := mockServ.events.Mock.Assert(
				mock.Call("Publish", mock.NotNil, mock.NotNil).Return(nil),
			); msg != "" {
				t.Errorf("test %d: %s", i, msg)
			}
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
	user1 := newMockUser("")
	user1.Validated = true
	user2 := newMockUser("")
	user2.Validated = true
	user2.Username = "admin"
	user2.Email = "admin@admin.com"
	user2.Name = "Admin"
	user2.Lastname = "Admin"
	user2.Roles = []Role{ADMIN}

	t.Run("Error", func(t *testing.T) {
		assert := assert.New(t)

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

			assert.True(errors.Compare(test.err, err), i)
		}
	})

	t.Run("OK", func(t *testing.T) {
		assert := assert.New(t)

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

			assert.Nil(err)

			if msg := mockServ.events.Mock.Assert(
				mock.Call("Publish", mock.NotNil, mock.NotNil).Return(mock.Nil),
			); msg != "" {
				t.Errorf("test %d: %s", i, msg)
			}
			userEvent, ok1 := mockServ.events.Mock.Calls[0].Args[0].(*UserEvent)
			opts, ok2 := mockServ.events.Mock.Calls[0].Args[1].(*events.Options)
			if !ok1 || !ok2 {
				t.Errorf("test %d: invalid conversion", i)
				continue
			}
			if userEvent.Type != "UserDeleted" {
				t.Errorf("test %d: invalid event type %s", i, userEvent.Type)
			}
			if userEvent.User.ID.Hex() != test.id {
				t.Errorf("test %d: invalid user id %s", i, userEvent.User.ID.Hex())
			}
			if opts.Exchange != "user" || opts.Route != "user.deleted" || opts.Queue != "" {
				t.Errorf("test %d: invalid options %#v", i, opts)
			}
		}
	})

}

func TestLogin(t *testing.T) {
	user1 := newMockUser(userID)
	user2 := newMockUser(adminID)
	user2.Username = "admin"
	user2.SetPassword("admin1234")
	user2.Email = "admin@admin.com"
	user3 := newMockUser("")
	user3.Username = "other-user"
	user3.Email = "other@user.com"

	t.Run("Error", func(t *testing.T) {
		assert := assert.New(t)

		req1 := &LoginRequest{
			Username: new(string),
			Password: new(string),
		}
		*req1.Username = "invalid-user"
		*req1.Password = "123456789"
		req2 := &LoginRequest{
			Username: new(string),
			Password: new(string),
		}
		*req2.Username = "user"
		*req2.Password = "invalid-password"
		req3 := &LoginRequest{
			Username: new(string),
			Password: new(string),
		}
		*req3.Username = "invalid-user"
		*req3.Password = "invalid-password"
		req4 := &LoginRequest{
			Email:    new(string),
			Password: new(string),
		}
		*req4.Email = "invalid@email.com"
		*req4.Password = "123456789"
		req5 := &LoginRequest{
			Email:    new(string),
			Password: new(string),
		}
		*req5.Email = "user@user.com"
		*req5.Password = "invalid-password"
		req6 := &LoginRequest{
			Email:    new(string),
			Password: new(string),
		}
		*req6.Email = "invalid@email.com"
		*req6.Password = "invalid-password"
		req7 := &LoginRequest{
			Username: new(string),
			Password: new(string),
		}
		*req7.Username = "other-user"
		*req7.Password = "123456789"
		allEmptyReq := &LoginRequest{}
		emptyPasswordReq := &LoginRequest{
			Username: new(string),
			Email:    new(string),
		}
		*emptyPasswordReq.Username = "user"
		*emptyPasswordReq.Email = "user@user.com"
		onlyPasswordReq := &LoginRequest{
			Password: new(string),
		}
		*onlyPasswordReq.Password = "123456789"

		tests := []struct {
			req *LoginRequest
			err error
		}{{
			req1,
			ErrInvalidUser,
		}, {
			req2,
			ErrInvalidUser,
		}, {
			req3,
			ErrInvalidUser,
		}, {
			req4,
			ErrInvalidUser,
		}, {
			req5,
			ErrInvalidUser,
		}, {
			req6,
			ErrInvalidUser,
		}, {
			req7,
			ErrInvalidUser,
		}, {
			allEmptyReq,
			ErrInvalidLogin,
		}, {
			emptyPasswordReq,
			ErrInvalidLogin,
		}, {
			onlyPasswordReq,
			ErrInvalidLogin,
		}}

		for i, test := range tests {
			serv := newMockService()
			serv.repo.populate(user1, user2)
			serv.authServ.populate(user1.ID.Hex(), user2.ID.Hex(), user3.ID.Hex())

			token, err := serv.Login(test.req)
			if assert.NotNil(err, i) {
				assert.True(errors.Compare(test.err, err), i)
			}
			assert.Nil(token, i)

			if test.req.Password == nil {
				continue
			}

			repoCalls := mock.Calls{}
			if test.req.Username != nil {
				repoCalls = append(repoCalls, mock.Call("FindByUsername", *test.req.Username))
			} else if test.req.Email != nil {
				repoCalls = append(repoCalls, mock.Call("FindByEmail", *test.req.Email))
			}

			if msg := serv.repo.Mock.Assert(repoCalls...); msg != "" {
				t.Errorf("test %d: %s", i, msg)
			}
		}
	})

	t.Run("OK", func(t *testing.T) {
		assert := assert.New(t)

		req1 := &LoginRequest{
			Username: new(string),
			Password: new(string),
		}
		*req1.Username = "user"
		*req1.Password = "123456789"
		req2 := &LoginRequest{
			Email:    new(string),
			Password: new(string),
		}
		*req2.Email = "user@user.com"
		*req2.Password = "123456789"
		req3 := &LoginRequest{
			Username: new(string),
			Email:    new(string),
			Password: new(string),
		}
		*req3.Username = "user"
		*req3.Email = "user@user.com"
		*req3.Password = "123456789"
		req4 := &LoginRequest{
			Username: new(string),
			Email:    new(string),
			Password: new(string),
		}
		*req4.Username = "admin"
		*req4.Email = "admin@admin.com"
		*req4.Password = "admin1234"
		req5 := &LoginRequest{
			Username: new(string),
			Password: new(string),
		}
		*req5.Username = "admin"
		*req5.Password = "admin1234"

		tests := []struct {
			req    *LoginRequest
			userID string
		}{{
			req1,
			user1.ID.Hex(),
		}, {
			req2,
			user1.ID.Hex(),
		}, {
			req3,
			user1.ID.Hex(),
		}, {
			req4,
			user2.ID.Hex(),
		}, {
			req5,
			user2.ID.Hex(),
		}}

		for i, test := range tests {
			serv := newMockService()
			serv.repo.populate(user1, user2)
			serv.authServ.populate(user1.ID.Hex(), user2.ID.Hex())

			token, err := serv.Login(test.req)
			assert.Nil(err)
			if assert.NotNil(token) {
				assert.Equal(test.userID, token.UserID)
			}

			repoCalls := mock.Calls{}

			if test.req.Username != nil {
				repoCalls = append(repoCalls, mock.Call("FindByUsername", *test.req.Username).Return(mock.NotNil, mock.Nil))
			} else if test.req.Email != nil {
				repoCalls = append(repoCalls, mock.Call("FindByEmail", *test.req.Email).Return(mock.NotNil, mock.Nil))
			}

			if msg := serv.repo.Mock.Assert(repoCalls...); msg != "" {
				t.Errorf("test %d: %s", i, msg)
			}

			if msg := serv.authServ.Mock.Assert(
				mock.Call("Create", mock.NotNil).Return(mock.NotNil, mock.Nil),
			); msg != "" {
				t.Errorf("test %d: %s", i, msg)
			}
		}
	})
}

func TestLogout(t *testing.T) {
	user := newMockUser("")

	t.Run("Error", func(t *testing.T) {
		assert := assert.New(t)

		tests := []struct {
			tokenStr func(serv *mockService) string
			err      error
		}{{
			func(serv *mockService) string {
				return "invalidToken123"
			},
			ErrInvalidUser,
		}, {
			func(serv *mockService) string {
				serv.repo.populate(user)
				return "user"
			},
			ErrInvalidUser,
		}, {
			func(serv *mockService) string {
				serv.authServ.populate(user.ID.Hex())
				return serv.authServ.tokensStr[user.ID.Hex()]
			},
			ErrInvalidUser,
		}}

		for i, test := range tests {
			serv := newMockService()
			tokenStr := test.tokenStr(serv)
			token := serv.authServ.tokens[tokenStr]

			err := serv.Logout(tokenStr)
			if assert.NotNil(err, i) {
				assert.True(errors.Compare(test.err, err), i)
			}
			assert.Len(serv.authServ.tokensStr, 0, i)

			if token != nil {
				if msg := serv.authServ.Mock.Assert(
					mock.Call("Validate", tokenStr).Return(token, mock.Nil),
					mock.Call("Invalidate", tokenStr).Return(token, mock.Nil),
				); msg != "" {
					t.Errorf("test %d: %s", i, msg)
				}

			} else {
				if msg := serv.authServ.Mock.Assert(
					mock.Call("Validate", tokenStr).Return(mock.Nil, mock.NotNil),
					mock.Call("Invalidate", tokenStr).Return(mock.Nil, mock.NotNil),
				); msg != "" {
					t.Errorf("test %d: %s", i, msg)
				}
			}

		}
	})

	t.Run("OK", func(t *testing.T) {
		assert := assert.New(t)

		serv := newMockService()
		serv.repo.populate(user)
		serv.authServ.populate(user.ID.Hex())
		tokenStr := serv.authServ.tokensStr[user.ID.Hex()]
		token := serv.authServ.tokens[tokenStr]

		err := serv.Logout(tokenStr)
		assert.Nil(err)

		if msg := serv.authServ.Mock.Assert(
			mock.Call("Validate", tokenStr).Return(token, mock.Nil),
			mock.Call("Invalidate", tokenStr).Return(token, mock.Nil),
		); msg != "" {
			t.Errorf("%s", msg)
		}
	})
}
