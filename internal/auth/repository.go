package auth

import (
	"encoding/json"

	"github.com/aboglioli/big-brother/pkg/cache"
	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/models"
)

// Errors
var (
	ErrRepositoryNotFound = errors.Internal.New("auth.repository.not_found")
	ErrRepositoryInsert   = errors.Internal.New("auth.repository.insert")
	ErrRepositoryDelete   = errors.Internal.New("auth.repository.delete")
)

// Interfaces
type Repository interface {
	FindByID(tokenID string) (*models.Token, error)
	Insert(token *models.Token) error
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

func (r *repositoryImpl) FindByID(tokenID string) (*models.Token, error) {
	v, err := r.cache.Get(tokenID)
	if v == nil || err != nil {
		return nil, ErrRepositoryNotFound.Wrap(err)
	}

	b, ok := v.([]byte)
	if !ok {
		return nil, ErrRepositoryNotFound.M("wrong conversion")
	}

	token := &models.Token{}
	if err := json.Unmarshal(b, token); err != nil {
		return nil, ErrRepositoryNotFound.Wrap(err)
	}

	return token, nil
}

func (r *repositoryImpl) Insert(token *models.Token) error {
	b, err := json.Marshal(token)
	if err != nil {
		return ErrRepositoryInsert.Wrap(err)
	}

	if err := r.cache.Set(token.ID, b, cache.NoExpiration); err != nil {
		return ErrRepositoryInsert.Wrap(err)
	}
	return nil
}

func (r *repositoryImpl) Delete(tokenID string) error {
	err := r.cache.Delete(tokenID)
	if err != nil {
		return ErrRepositoryDelete.Wrap(err)
	}
	return nil
}
