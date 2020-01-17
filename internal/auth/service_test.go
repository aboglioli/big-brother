package auth

import (
	"testing"

	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateToken(t *testing.T) {
	mTokenStr := "encoded.token"

	tests := []struct {
		name   string
		userID string
		err    error
		mock   func(*mockService)
	}{{
		"empty userID",
		"",
		ErrCreate,
		nil,
	}, {
		"repo error",
		"user123",
		ErrCreate.Wrap(ErrRepositoryInsert),
		func(s *mockService) {
			s.enc.On("Encode", mock.Anything).Return(mTokenStr, nil)
			s.repo.On("Insert", mock.AnythingOfType("*models.Token")).Return(ErrRepositoryInsert)
		},
	}, {
		"user123",
		"user123",
		nil,
		func(s *mockService) {
			s.enc.On("Encode", mock.Anything).Return(mTokenStr, nil)
			s.repo.On("Insert", mock.AnythingOfType("*models.Token")).Return(nil)
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			serv := newMockService()
			if test.mock != nil {
				test.mock(serv)
			}

			tokenStr, err := serv.Create(test.userID)

			if test.err != nil { // Error
				if assert.NotNil(err) {
					errors.Assert(t, test.err, err)
				}
				assert.Empty(tokenStr)
			} else { // OK
				assert.Nil(err)
				if assert.NotEmpty(tokenStr) {
					assert.Equal(mTokenStr, tokenStr)
				}
				t, ok := serv.repo.Calls[0].Arguments[0].(*models.Token)
				assert.True(ok)
				assert.Equal(test.userID, t.UserID)
			}
			serv.enc.AssertExpectations(t)
			serv.repo.AssertExpectations(t)
		})
	}
}

func TestValidate(t *testing.T) {
	mToken := models.NewToken("user123")
	mTokenStr := "encoded.token"

	tests := []struct {
		name     string
		tokenStr string
		err      error
		mock     func(s *mockService)
	}{{
		"empty tokenStr",
		"",
		ErrValidate.Wrap(ErrTokenDecode),
		func(s *mockService) {
			s.enc.On("Decode", "").Return("", ErrTokenDecode)
		},
	}, {
		"invalid token",
		"eyJ.eyJj.fxg",
		ErrValidate.Wrap(ErrTokenDecode),
		func(s *mockService) {
			s.enc.On("Decode", "eyJ.eyJj.fxg").Return("", ErrTokenDecode)
		},
	}, {
		"not saved and valid token",
		mTokenStr,
		ErrValidate.Wrap(ErrRepositoryNotFound),
		func(s *mockService) {
			s.enc.On("Decode", mTokenStr).Return("token123", nil)
			s.repo.On("FindByID", "token123").Return(nil, ErrRepositoryNotFound)
		},
	}, {
		"valid token",
		mTokenStr,
		nil,
		func(s *mockService) {
			s.enc.On("Decode", mTokenStr).Return(mToken.ID, nil)
			s.repo.On("FindByID", mToken.ID).Return(mToken, nil)
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			serv := newMockService()
			if test.mock != nil {
				test.mock(serv)
			}

			token, err := serv.Validate(test.tokenStr)

			if test.err != nil { // Error
				if assert.NotNil(err) {
					errors.Assert(t, test.err, err)
				}
				assert.Nil(token)
			} else { // OK
				assert.Nil(err)
				if assert.NotNil(token) {
					assert.Equal(mToken.ID, token.ID)
					assert.Equal(mToken.UserID, token.UserID)
				}
			}
			serv.enc.AssertExpectations(t)
			serv.repo.AssertExpectations(t)
		})
	}
}

func TestInvalidate(t *testing.T) {
	mToken := models.NewToken("user123")
	mTokenStr := "encoded.token"

	tests := []struct {
		name     string
		tokenStr string
		err      error
		mock     func(s *mockService)
	}{{
		"empty tokenStr",
		"",
		ErrInvalidate.Wrap(ErrValidate.Wrap(ErrTokenDecode)),
		func(s *mockService) {
			s.enc.On("Decode", "").Return("", ErrTokenDecode)
		},
	}, {
		"invalid tokenStr",
		"asd.zxc.ey88",
		ErrInvalidate.Wrap(ErrValidate.Wrap(ErrTokenDecode)),
		func(s *mockService) {
			s.enc.On("Decode", "asd.zxc.ey88").Return("", ErrTokenDecode)
		},
	}, {
		"not saved and valid tokenStr",
		mTokenStr,
		ErrInvalidate.Wrap(ErrValidate.Wrap(ErrRepositoryNotFound)),
		func(s *mockService) {
			s.enc.On("Decode", mTokenStr).Return(mToken.ID, nil)
			s.repo.On("FindByID", mToken.ID).Return(nil, ErrRepositoryNotFound)
		},
	}, {
		"error on deleting",
		mTokenStr,
		ErrInvalidate.Wrap(ErrRepositoryDelete),
		func(s *mockService) {
			s.enc.On("Decode", mTokenStr).Return(mToken.ID, nil)
			s.repo.On("FindByID", mToken.ID).Return(mToken, nil)
			s.repo.On("Delete", mToken.ID).Return(ErrRepositoryDelete)
		},
	}, {
		"valid and saved tokenStr",
		mTokenStr,
		nil,
		func(s *mockService) {
			s.enc.On("Decode", mTokenStr).Return(mToken.ID, nil)
			s.repo.On("FindByID", mToken.ID).Return(mToken, nil)
			s.repo.On("Delete", mToken.ID).Return(nil)
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			serv := newMockService()
			if test.mock != nil {
				test.mock(serv)
			}

			token, err := serv.Invalidate(test.tokenStr)

			if test.err != nil { // Error
				if assert.NotNil(err) {
					errors.Assert(t, test.err, err)
				}
				assert.Nil(token)
			} else { // OK
				assert.Nil(err)
				if assert.NotNil(token) {
					assert.Equal(mToken.ID, token.ID)
					assert.Equal(mToken.UserID, token.UserID)
				}
			}
			serv.enc.AssertExpectations(t)
			serv.repo.AssertExpectations(t)
		})
	}
}
