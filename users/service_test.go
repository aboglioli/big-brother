package users

import (
	"reflect"
	"testing"

	"github.com/aboglioli/big-brother/errors"
	"github.com/aboglioli/big-brother/tools/mock"
)

func TestGetByID(t *testing.T) {
	user1 := newMockUser()
	user2 := newMockUser()
	user2.Validated = true
	user3 := newMockUser()
	user3.Validated = true
	user3.Active = false

	tests := []struct {
		in   string
		out  error
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
		u, err := mockServ.GetByID(test.in)

		if err != nil {
			if test.out != nil {
				expectedErr := test.out.(errors.Error)
				err := err.(errors.Error)
				if expectedErr.Code != err.Code {
					t.Errorf("test %d:\n-expected:%#v\n-actual:  %#v", i, test.out, err)
				}

			} else {
				t.Errorf("test %d:\n-expected: nil error\n-actual:%#v", i, err)
			}
		}

		if err == nil && !reflect.DeepEqual(u, test.user) {
			t.Errorf("test %d:\n-expected:%#v\n-actual:  %#v", i, test.user, u)
		}

		mockServ.repo.Mock.Assert(t,
			mock.Call("FindByID", test.in),
		)
	}
}

func TestCreate(t *testing.T) {
	user := newMockUser()
	user.Validated = true

	t.Run("Error", func(t *testing.T) {
		tests := []struct {
			in               *CreateRequest
			existingUsername bool
			existingEmail    bool
		}{{
			&CreateRequest{"admin", "1234567", "admin@admin.com"},
			false,
			false,
		}, {
			&CreateRequest{"user", "123456789", "admin@admi.com"},
			true,
			false,
		}, {
			&CreateRequest{"admin", "12345678", "user@user.com"},
			false,
			true,
		}}

		for i, test := range tests {
			mockServ := newMockService()
			mockServ.repo.populate(newMockUser())
			u, err := mockServ.Create(test.in)

			if u != nil {
				t.Errorf("test %d: expected nil user", i)
			}

			if err == nil {
				t.Errorf("test %d: expected error", i)
			}

			call1 := mock.Call("FindByUsername", test.in.Username)
			call2 := mock.Call("FindByEmail", test.in.Email)

			if test.existingUsername {
				call1 = call1.Return(mock.NotNil, mock.Nil)
			} else {
				call1 = call1.Return(mock.Nil, mock.NotNil)
			}

			if test.existingEmail {
				call2 = call2.Return(mock.NotNil, mock.Nil)
			} else {
				call2 = call2.Return(mock.Nil, mock.NotNil)
			}

			mockServ.repo.Mock.Assert(t,
				call1,
				call2,
			)

			mockServ.validator.Mock.Assert(t,
				mock.Call("ValidatePassword", test.in.Password),
				mock.Call("ValidateSchema", mock.NotNil),
			)
		}
	})

	t.Run("OK", func(t *testing.T) {
		tests := []struct {
			in  *CreateRequest
			out *User
		}{{
			&CreateRequest{"admin", "123456789", "admin@admin.com"},
			&User{
				Username: "admin",
				Email:    "admin@admin",
			},
		}}

		for i, test := range tests {
			mockServ := newMockService()
			user, err := mockServ.Create(test.in)

			if err != nil {
				t.Errorf("test %d: err not expected", i)
				break
			}

			if user == nil {
				t.Errorf("test %d: user expected", i)
			}

			if user.Username != test.out.Username || user.Email != user.Email {
				t.Errorf("test %d:\n-expected:%#v\n-actual:  %#v", i, test.out, user)
			}

			if user.Password == test.in.Password || len(user.Password) < 10 {
				t.Errorf("test %d: password wrong hashing: %s", i, user.Password)
			}

			mockServ.repo.Mock.Assert(t,
				mock.Call("FindByUsername", test.in.Username).Return(mock.Nil, mock.NotNil),
				mock.Call("FindByEmail", test.in.Email).Return(mock.Nil, mock.NotNil),
				mock.Call("Insert", mock.NotNil).Return(nil),
			)
			insertedUser := mockServ.repo.Mock.Calls[2].Args[0]

			if !reflect.DeepEqual(user, insertedUser) {
				t.Errorf("test %d: inserted user not equal returned user\n-expected:%#v\n-actual:  %#v", i, user, insertedUser)
			}
		}
	})

}
