package auth

import (
	"github.com/aboglioli/big-brother/pkg/cache"
	"github.com/aboglioli/big-brother/pkg/errors"
)

// Errors
var (
	ErrRepositoryNotFound = errors.Internal.New("auth.repository.not_found")
	ErrRepositoryInsert   = errors.Internal.New("auth.repository.insert")
	ErrRepositoryDelete   = errors.Internal.New("auth.repository.delete")
)

// Interfaces
type Repository interface {
	FindByID(tokenID string) (*Token, error)
	Insert(token *Token) error
	Delete(tokenID string) error
}

// Implementations
type repositoryImpl struct {
	cache cache.Cache
}

func NewRepository(cache cache.Cache) Repository {
	return &repositoryImpl{
		cache: cache,
	}
}

func (r *repositoryImpl) FindByID(tokenID string) (*Token, error) {
	token := r.cache.Get(tokenID)
	if token == nil {
		return nil, ErrRepositoryNotFound
	}
	if t, ok := token.(*Token); ok {
		return t, nil
	}

	return nil, ErrRepositoryNotFound
}

func (r *repositoryImpl) Insert(token *Token) error {
	r.cache.Set(token.ID.Hex(), token, cache.NoExpiration)
	return nil
}

func (r *repositoryImpl) Delete(tokenID string) error {
	r.cache.Delete(tokenID)
	return nil
}
