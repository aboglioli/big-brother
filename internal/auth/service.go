package auth

import (
	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/models"
)

// Errors
var (
	ErrCreate     = errors.Status.New("auth.service.create")
	ErrValidate   = errors.Status.New("auth.service.validate")
	ErrInvalidate = errors.Status.New("auth.service.invalidate")
)

// Interface
type Service interface {
	Create(userID string) (string, error)
	Validate(tokenStr string) (*models.Token, error)
	Invalidate(tokenStr string) (*models.Token, error)
}

// Implementation
type service struct {
	repo Repository
	enc  Encoder
}

func NewService(repo Repository) Service {
	enc := NewEncoder()
	return &service{
		repo: repo,
		enc:  enc,
	}
}

func (s *service) Create(userID string) (string, error) {
	if userID == "" {
		return "", ErrCreate
	}

	token := models.NewToken(userID)
	tokenStr, err := s.enc.Encode(token.ID)
	if tokenStr == "" || err != nil {
		return "", ErrCreate.Wrap(err)
	}

	if err := s.repo.Insert(token); err != nil {
		return "", ErrCreate.Wrap(err)
	}

	return tokenStr, nil
}

func (s *service) Validate(tokenStr string) (*models.Token, error) {
	tokenID, err := s.enc.Decode(tokenStr)
	if tokenID == "" || err != nil {
		return nil, ErrValidate.Wrap(err)
	}

	token, err := s.repo.FindByID(tokenID)
	if token == nil || err != nil {
		return nil, ErrValidate.Wrap(err)
	}

	return token, nil
}

func (s *service) Invalidate(tokenStr string) (*models.Token, error) {
	token, err := s.Validate(tokenStr)
	if err != nil {
		return nil, ErrInvalidate.Wrap(err)
	}

	if err := s.repo.Delete(token.ID); err != nil {
		return nil, ErrInvalidate.Wrap(err)
	}

	return token, nil
}
