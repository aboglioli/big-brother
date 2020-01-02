package auth

import (
	"fmt"

	"github.com/aboglioli/big-brother/cache"
	"github.com/aboglioli/big-brother/errors"
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
	token := r.cache.Get(fmt.Sprintf("auth:%s", tokenID))
	if t, ok := token.(*Token); ok {
		return t, nil
	}

	return nil, ErrRepositoryNotFound
}

func (r *repositoryImpl) Insert(token *Token) error {
	r.cache.Set(fmt.Sprintf("auth:%s", token.ID.Hex()), token, cache.NoExpiration)
	return nil
}

func (r *repositoryImpl) Delete(tokenID string) error {
	r.cache.Delete(fmt.Sprintf("auth:%s", tokenID))
	return nil
}
