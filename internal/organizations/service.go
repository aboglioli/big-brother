package organizations

import (
	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/models"
)

// Interface
type Service interface {
	GetByID(id string) (*models.Organization, error)
	Search(query string) ([]*models.Organization, error)

	Create(req *CreateRequest) (*models.Organization, error)
	Update(req *UpdateRequest) (*models.Organization, error)
	Delete(id string) error
}

// Request DTOs
type CreateRequest struct {
	Name string
}

type UpdateRequest struct {
	Name string
}

// Implementation
type service struct {
	repo      Repository
	validator Validator
}

func NewService(repo Repository) Service {
	return &service{
		repo:      repo,
		validator: NewValidator(),
	}
}

func (s *service) GetByID(id string) (*models.Organization, error) {
	return s.getByID(id)
}

func (s *service) Search(query string) ([]*models.Organization, error) {
	if query == "" {
		return []*models.Organization{}, nil
	}
	orgs, err := s.repo.SearchByName(query)
	if err != nil {
		return nil, errors.ErrNotFound.Wrap(err)
	}
	return orgs, nil
}

func (s *service) Create(req *CreateRequest) (*models.Organization, error) {
	return nil, nil
}

func (s *service) Update(req *UpdateRequest) (*models.Organization, error) {
	return nil, nil
}

func (s *service) Delete(id string) error {
	return nil
}

func (s *service) getByID(id string) (*models.Organization, error) {
	if id == "" {
		return nil, errors.ErrNotFound.M("invalid id")
	}

	org, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.ErrNotFound.Wrap(err)
	}

	if err := s.validator.Status(org); err != nil {
		return nil, err
	}

	return org, nil
}
