package auth

import (
	"github.com/aboglioli/big-brother/errors"
	"github.com/aboglioli/big-brother/users"
)

// Errors
var (
	ErrCreate       = errors.Status.New("auth.service.create")
	ErrUnauthorized = errors.Status.New("auth.service.unauthorized")
)

// Interfaces
type Service interface {
	Create(user *users.User) (*Token, error)
	Validate(tokenStr string) (*Token, error)
	Invalidate(tokenStr string) error
}

// Implementations
type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{
		repo: repo,
	}
}

func (s *serviceImpl) Create(user *users.User) (*Token, error) {
	token := NewToken(user.ID.Hex())
	if err := s.repo.Insert(token); err != nil {
		return nil, ErrCreate.Wrap(err)
	}

	return token, nil
}

func (s *serviceImpl) Validate(tokenStr string) (*Token, error) {
	token, err := decodeToken(tokenStr)
	if err != nil {
		return nil, ErrUnauthorized.Wrap(err)
	}

	t, err := s.repo.FindByID(token.ID.Hex())
	if t == nil || err != nil {
		return nil, ErrUnauthorized.Wrap(err)
	}

	return token, nil
}

func (s *serviceImpl) Invalidate(tokenStr string) error {
	token, err := s.Validate(tokenStr)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(token.ID.Hex()); err != nil {
		return ErrUnauthorized.Wrap(err)
	}

	return nil
}
