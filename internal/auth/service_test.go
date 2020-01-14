package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateToken(t *testing.T) {
	assert := assert.New(t)

	serv := newMockService()
	userID := "user123"
	token := NewToken(userID)
	serv.repo.On("Insert", mock.Anything).Return(nil)

	token, err := serv.Create(userID)
	assert.Nil(err)
	if assert.NotNil(token) {
		assert.NotEmpty(token.ID.Hex())
		assert.Equal(token.UserID, userID)
	}

	serv.repo.AssertExpectations(t)
}

func TestValidate(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		assert := assert.New(t)

		invalid := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVkQXQiOiIyMDIwLTAxLTA1VDAxOjA3OjQxLjcxMTQ0ODY1Ni0wMzowMCIsImlkIjoiNWUxMTYxMGQzNTA1MjI0YTlmNzJmN2Q2IiwidXNlcklkIjoiMTIzNCJ9.fxg_UZMR8fBaVluRmekEslf453DlJ_oA_QX8fv3QkFQ"

		tests := []struct {
			tokenStr string
			fn       func(*mockRepository)
		}{{
			"",
			nil,
		}, {
			"123",
			nil,
		}, {
			"456789abc",
			nil,
		}, {
			invalid,
			func(repo *mockRepository) {
				repo.On("FindByID", "5e11610d3505224a9f72f7d6").Return(&Token{}, ErrUnauthorized)
			},
		}, {
			"123456789",
			nil,
		}}

		for i, test := range tests {
			serv := newMockService()
			if test.fn != nil {
				test.fn(serv.repo)
			}
			token, err := serv.Validate(test.tokenStr)
			assert.NotNil(err, i)
			assert.Nil(token, i)
			serv.repo.AssertExpectations(t)
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
			tokenStr, err := token.Encode()
			assert.Nil(err, i)

			serv.repo.On("FindByID", token.ID.Hex()).Return(token, nil)

			validatedToken, err := serv.Validate(tokenStr)
			assert.Nil(err)
			if assert.NotNil(validatedToken) {
				assert.Equal(validatedToken.ID.Hex(), token.ID.Hex())
				assert.Equal(validatedToken.UserID, token.UserID)
			}

			serv.repo.AssertExpectations(t)
		}
	})
}

func TestInvalidate(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		assert := assert.New(t)

		tests := []struct {
			tokenStr string
		}{{""}, {"1234"}, {"123456"}, {"abc123"}}

		for i, test := range tests {
			serv := newMockService()
			token, err := serv.Invalidate(test.tokenStr)
			assert.NotNil(err, i)
			assert.Nil(token)
		}
	})

	t.Run("OK", func(t *testing.T) {
		assert := assert.New(t)

		userID := "user123"
		tests := []struct {
			userID string
		}{{userID}, {"123456"}, {"abc123"}}

		for i, test := range tests {
			serv := newMockService()
			token := NewToken(test.userID)
			tokenStr, err := token.Encode()
			assert.Nil(err, i)
			if err != nil {
				t.Errorf("test %d: error not expected: %s", i, err)
			}

			serv.repo.On("FindByID", token.ID.Hex()).Return(&Token{
				ID:        token.ID,
				UserID:    token.UserID,
				CreatedAt: token.CreatedAt,
			}, nil)
			serv.repo.On("Delete", token.ID.Hex()).Return(nil)

			invalidatedToken, err := serv.Invalidate(tokenStr)
			assert.Nil(err, i)
			assert.Equal(token.ID.Hex(), invalidatedToken.ID.Hex(), i)

			serv.repo.AssertExpectations(t)
		}
	})
}
