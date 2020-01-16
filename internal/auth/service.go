package auth

import (
	"github.com/aboglioli/big-brother/pkg/errors"
)

// Errors
var (
	ErrCreate     = errors.Status.New("auth.service.create")
	ErrValidate   = errors.Status.New("auth.service.validate")
	ErrInvalidate = errors.Status.New("auth.service.invalidate")
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
		return nil, ErrValidate.Wrap(err)
	}

	t, err := s.repo.FindByID(token.ID)
	if t == nil || err != nil {
		return nil, ErrValidate.Wrap(err)
	}

	return t, nil
}

func (s *serviceImpl) Invalidate(tokenStr string) (*Token, error) {
	token, err := s.Validate(tokenStr)
	if err != nil {
		return nil, ErrInvalidate.Wrap(err)
	}

	if err := s.repo.Delete(token.ID); err != nil {
		return nil, ErrInvalidate.Wrap(err)
	}

	return token, nil
}
