package auth

import (
	"testing"

	"github.com/aboglioli/big-brother/tools/mock"
	"github.com/aboglioli/big-brother/users"
)

func TestCreateToken(t *testing.T) {
	serv := newMockService()
	user := users.NewUser()

	token, err := serv.Create(user)
	if token == nil || err != nil {
		t.Errorf("expected token, got error %#v", err)
		return
	}
	if len(token.ID.Hex()) < 6 || token.UserID != user.ID.Hex() {
		t.Errorf("invaid token %#v", token)

	}
	serv.repo.Mock.Assert(t,
		mock.Call("Insert", token).Return(mock.Nil),
	)

	rawSavedToken := serv.repo.Repo.cache.Get(token.ID.Hex())
	if rawSavedToken == nil {
		t.Errorf("token is not saved in repository")
		return

	}
	savedToken, ok := rawSavedToken.(*Token)
	if !ok {
		t.Errorf("invalid conversion from repository")
		return
	}
	if savedToken.ID.Hex() != token.ID.Hex() || savedToken.UserID != token.UserID {
		t.Errorf("savedToken\n-expected:%#v\n-actual:  %#v", token, savedToken)
	}
}

func TestValidate(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		invalid := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVkQXQiOiIyMDIwLTAxLTA1VDAxOjA3OjQxLjcxMTQ0ODY1Ni0wMzowMCIsImlkIjoiNWUxMTYxMGQzNTA1MjI0YTlmNzJmN2Q2IiwidXNlcklkIjoiMTIzNCJ9.fxg_UZMR8fBaVluRmekEslf453DlJ_oA_QX8fv3QkFQ"

		tests := []struct {
			tokenStr string
		}{{""}, {"123"}, {"456789abc"}, {invalid}, {"123456789"}}

		for i, test := range tests {
			serv := newMockService()
			token, err := serv.Validate(test.tokenStr)
			if token != nil || err == nil {
				t.Errorf("test %d: error expected, got token", i)
				continue
			}
		}
	})

	t.Run("OK", func(t *testing.T) {
		user := users.NewUser()

		tests := []struct {
			userID string
		}{{user.ID.Hex()}, {"123456"}}

		for i, test := range tests {
			serv := newMockService()
			token := NewToken(test.userID)
			tokenStr, err := token.Encode()
			if err != nil {
				t.Errorf("test %d: errors not expected: %s", i, err)
				continue
			}
			serv.repo.populate(token)

			validatedToken, err := serv.Validate(tokenStr)
			if validatedToken == nil || err != nil {
				t.Errorf("test %d: validated token expected, got error: %s", i, err)
			}
			if validatedToken.ID.Hex() != token.ID.Hex() || validatedToken.UserID != token.UserID {
				t.Errorf("test %d: tokens are not equal\n-expected:%#v\n-actual:  %#v", i, token, validatedToken)
			}

			serv.repo.Mock.Assert(t,
				mock.Call("FindByID", token.ID.Hex()).Return(token, mock.Nil),
			)
		}
	})
}

func TestInvalidate(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		tests := []struct {
			tokenStr string
		}{{""}, {"1234"}, {"123456"}, {"abc123"}}

		for i, test := range tests {
			serv := newMockService()
			err := serv.Invalidate(test.tokenStr)
			if err == nil {
				t.Errorf("test %d: expected error", i)
			}
		}
	})

	t.Run("OK", func(t *testing.T) {
		user := users.NewUser()
		tests := []struct {
			userID string
		}{{user.ID.Hex()}, {"123456"}, {"abc123"}}

		for i, test := range tests {
			serv := newMockService()
			token := NewToken(test.userID)
			tokenStr, err := token.Encode()
			if err != nil {
				t.Errorf("test %d: error not expected: %s", i, err)
			}
			serv.repo.populate(token)

			err = serv.Invalidate(tokenStr)
			if err != nil {
				t.Errorf("test %d: error not expected: %s", i, err)
			}

			serv.repo.Mock.Assert(t,
				mock.Call("FindByID", mock.NotNil).Return(mock.NotNil, mock.Nil),
				mock.Call("Delete", token.ID.Hex()).Return(mock.Nil),
			)
		}
	})
}
