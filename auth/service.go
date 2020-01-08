package auth

import (
	"github.com/aboglioli/big-brother/errors"
)

// Errors
var (
	ErrCreate       = errors.Status.New("auth.service.create")
	ErrUnauthorized = errors.Status.New("auth.service.unauthorized")
)

// Interfaces
type Service interface {
	Create(userID string) (*Token, error)
	Validate(tokenStr string) (*Token, error)
	Invalidate(tokenStr string) (*Token, error)
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

func (s *serviceImpl) Create(userID string) (*Token, error) {
	token := NewToken(userID)
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

func (s *serviceImpl) Invalidate(tokenStr string) (*Token, error) {
	token, err := s.Validate(tokenStr)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Delete(token.ID.Hex()); err != nil {
		return nil, ErrUnauthorized.Wrap(err)
	}

	return token, nil
}
