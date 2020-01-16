package auth

import (
	"testing"

	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateToken(t *testing.T) {
	tests := []struct {
		name   string
		userID string
		err    error
		mock   func(*mockService)
	}{{
		"empty userID",
		"",
		nil,
		func(s *mockService) {
			s.repo.On("Insert", mock.AnythingOfType("*auth.Token")).Return(nil)
		},
	}, {
		"user123",
		"user123",
		nil,
		func(s *mockService) {
			s.repo.On("Insert", mock.AnythingOfType("*auth.Token")).Return(nil)
		},
	}, {
		"repo error",
		"user123",
		ErrCreate.Wrap(ErrRepositoryInsert),
		func(s *mockService) {
			s.repo.On("Insert", mock.AnythingOfType("*auth.Token")).Return(ErrRepositoryInsert)
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			serv := newMockService()
			if test.mock != nil {
				test.mock(serv)
			}

			token, err := serv.Create(test.userID)

			if test.err != nil { // Error
				if assert.NotNil(err) {
					errors.Assert(t, test.err, err)
				}
				assert.Nil(token)
			} else { // OK
				assert.Nil(err)
				if assert.NotNil(token) {
					assert.NotEmpty(token.ID)
					assert.Equal(test.userID, token.UserID)
				}
				serv.repo.AssertCalled(t, "Insert", token)
			}
			serv.repo.AssertExpectations(t)
		})
	}
}

func TestValidate(t *testing.T) {
	mToken := NewToken("user123")
	mTokenStr, err := mToken.Encode()
	require.Nil(t, err)

	tests := []struct {
		name     string
		tokenStr string
		err      error
		mock     func(s *mockService)
	}{{
		"empty tokenStr",
		"",
		ErrValidate.Wrap(ErrTokenDecode),
		nil,
	}, {
		"invalid token",
		"eyJ.eyJj.fxg",
		ErrValidate.Wrap(ErrTokenDecode),
		nil,
	}, {
		"not saved and valid token",
		mTokenStr,
		ErrValidate.Wrap(ErrRepositoryNotFound),
		func(s *mockService) {
			s.repo.On("FindByID", mToken.ID).Return(nil, ErrRepositoryNotFound)
		},
	}, {
		"valid token",
		mTokenStr,
		nil,
		func(s *mockService) {
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
			serv.repo.AssertExpectations(t)
		})
	}
}

func TestInvalidate(t *testing.T) {
	mToken := NewToken("user123")
	mTokenStr, err := mToken.Encode()
	require.Nil(t, err)

	tests := []struct {
		name     string
		tokenStr string
		err      error
		mock     func(s *mockService)
	}{{
		"empty tokenStr",
		"",
		ErrInvalidate.Wrap(ErrValidate.Wrap(ErrTokenDecode)),
		nil,
	}, {
		"invalid tokenStr",
		"asd.zxc.ey88",
		ErrInvalidate.Wrap(ErrValidate.Wrap(ErrTokenDecode)),
		nil,
	}, {
		"not saved and valid tokenStr",
		mTokenStr,
		ErrInvalidate.Wrap(ErrValidate.Wrap(ErrRepositoryNotFound)),
		func(s *mockService) {
			s.repo.On("FindByID", mToken.ID).Return(nil, ErrRepositoryNotFound)
		},
	}, {
		"error on deleting",
		mTokenStr,
		ErrInvalidate.Wrap(ErrRepositoryDelete),
		func(s *mockService) {
			s.repo.On("FindByID", mToken.ID).Return(mToken, nil)
			s.repo.On("Delete", mToken.ID).Return(ErrRepositoryDelete)
		},
	}, {
		"valid and saved tokenStr",
		mTokenStr,
		nil,
		func(s *mockService) {
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
			serv.repo.AssertExpectations(t)
		})
	}
}
