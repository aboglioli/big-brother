package auth

import (
	"testing"

	"github.com/aboglioli/big-brother/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func TestCreateToken(t *testing.T) {
	assert := assert.New(t)

	serv := newMockService()
	userID := "user123"

	token, err := serv.Create(userID)
	assert.Nil(err)
	if assert.NotNil(token) {
		assert.NotEmpty(token.ID.Hex())
		assert.Equal(token.UserID, userID)
	}

	if msg := serv.repo.Mock.Assert(
		mock.Call("Insert", token).Return(mock.Nil),
	); msg != "" {
		t.Errorf("%s", msg)
	}

	rawSavedToken := serv.repo.Repo.cache.Get(token.ID.Hex())
	assert.NotNil(rawSavedToken)

	savedToken, ok := rawSavedToken.(*Token)
	if assert.Equal(ok, true) {
		assert.Equal(savedToken.ID.Hex(), token.ID.Hex())
		assert.Equal(savedToken.UserID, token.UserID)
	}
}

func TestValidate(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		assert := assert.New(t)

		invalid := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVkQXQiOiIyMDIwLTAxLTA1VDAxOjA3OjQxLjcxMTQ0ODY1Ni0wMzowMCIsImlkIjoiNWUxMTYxMGQzNTA1MjI0YTlmNzJmN2Q2IiwidXNlcklkIjoiMTIzNCJ9.fxg_UZMR8fBaVluRmekEslf453DlJ_oA_QX8fv3QkFQ"

		tests := []struct {
			tokenStr string
		}{{""}, {"123"}, {"456789abc"}, {invalid}, {"123456789"}}

		for _, test := range tests {
			serv := newMockService()
			token, err := serv.Validate(test.tokenStr)
			assert.NotNil(err)
			assert.Nil(token)
		}
	})

	t.Run("OK", func(t *testing.T) {
		assert := assert.New(t)

		userID := "user123"

		tests := []struct {
			userID string
		}{{userID}, {"123456"}}

		for i, test := range tests {
			serv := newMockService()
			token := NewToken(test.userID)
			serv.repo.populate(token)

			tokenStr, err := token.Encode()
			assert.Nil(err)

			validatedToken, err := serv.Validate(tokenStr)
			assert.Nil(err)
			if assert.NotNil(validatedToken) {
				assert.Equal(validatedToken.ID.Hex(), token.ID.Hex())
				assert.Equal(validatedToken.UserID, token.UserID)
			}

			if msg := serv.repo.Mock.Assert(
				mock.Call("FindByID", token.ID.Hex()).Return(token, mock.Nil),
			); msg != "" {
				t.Errorf("test %d: %s", i, msg)
			}
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
			token, err := serv.Invalidate(test.tokenStr)
			if token != nil || err == nil {
				t.Errorf("test %d: expected error", i)
			}
		}
	})

	t.Run("OK", func(t *testing.T) {
		userID := "user123"
		tests := []struct {
			userID string
		}{{userID}, {"123456"}, {"abc123"}}

		for i, test := range tests {
			serv := newMockService()
			token := NewToken(test.userID)
			tokenStr, err := token.Encode()
			if err != nil {
				t.Errorf("test %d: error not expected: %s", i, err)
			}
			serv.repo.populate(token)

			invalidatedToken, err := serv.Invalidate(tokenStr)
			if err != nil {
				t.Errorf("test %d: error not expected: %s", i, err)
			}
			if token.ID.Hex() != invalidatedToken.ID.Hex() {
				t.Errorf("test %d: different tokens", i)
			}

			if msg := serv.repo.Mock.Assert(
				mock.Call("FindByID", mock.NotNil).Return(mock.NotNil, mock.Nil),
				mock.Call("Delete", token.ID.Hex()).Return(mock.Nil),
			); msg != "" {
				t.Errorf("test %d, %s", i, msg)
			}
		}
	})
}
