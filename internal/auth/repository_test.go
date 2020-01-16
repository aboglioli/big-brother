package auth

import (
	"encoding/json"
	"testing"

	"github.com/aboglioli/big-brother/mocks"
	"github.com/aboglioli/big-brother/pkg/cache"
	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindByID(t *testing.T) {
	mToken := models.NewToken("user123")
	mBytes, err := json.Marshal(mToken)
	require.Nil(t, err)

	tests := []struct {
		name string
		id   string
		err  error
		mock func(c *mocks.MockCache)
	}{{
		"empty id",
		"",
		ErrRepositoryNotFound,
		func(c *mocks.MockCache) {
			c.On("Get", "").Return(nil, cache.ErrCacheNotFound)
		},
	}, {
		"not found",
		"token123",
		ErrRepositoryNotFound,
		func(c *mocks.MockCache) {
			c.On("Get", "token123").Return(nil, cache.ErrCacheNotFound)
		},
	}, {
		"existing token",
		mToken.ID,
		nil,
		func(c *mocks.MockCache) {
			c.On("Get", mToken.ID).Return(mBytes, nil)
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			mCache := mocks.NewMockCache()
			repo := NewRepository(mCache)
			if test.mock != nil {
				test.mock(mCache)
			}

			token, err := repo.FindByID(test.id)

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
			mCache.AssertExpectations(t)
		})
	}
}

func TestInsert(t *testing.T) {
	mToken := models.NewToken("user123")
	mBytes, err := json.Marshal(mToken)
	require.Nil(t, err)

	tests := []struct {
		name  string
		token *models.Token
		err   error
		mock  func(c *mocks.MockCache)
	}{{
		"error on insert",
		mToken,
		ErrRepositoryInsert,
		func(c *mocks.MockCache) {
			c.On("Set", mToken.ID, mBytes).Return(cache.ErrCacheSet)
		},
	}, {
		"success",
		mToken,
		nil,
		func(c *mocks.MockCache) {
			c.On("Set", mToken.ID, mBytes).Return(nil)
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			mCache := mocks.NewMockCache()
			repo := NewRepository(mCache)
			if test.mock != nil {
				test.mock(mCache)
			}

			err := repo.Insert(test.token)

			if test.err != nil {
				if assert.NotNil(err) {
					errors.Assert(t, test.err, err)
				}
			} else {
				assert.Nil(err)
			}
			mCache.AssertExpectations(t)
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name string
		id   string
		err  error
		mock func(c *mocks.MockCache)
	}{{
		"not found",
		"token123",
		ErrRepositoryDelete,
		func(c *mocks.MockCache) {
			c.On("Delete", "token123").Return(cache.ErrCacheDelete)
		},
	}, {
		"success",
		"token123",
		nil,
		func(c *mocks.MockCache) {
			c.On("Delete", "token123").Return(nil)
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			mCache := mocks.NewMockCache()
			repo := NewRepository(mCache)
			if test.mock != nil {
				test.mock(mCache)
			}

			err := repo.Delete(test.id)

			if test.err != nil {
				if assert.NotNil(err) {
					errors.Assert(t, test.err, err)
				}
			} else {
				assert.Nil(err)
			}
			mCache.AssertExpectations(t)
		})
	}
}
