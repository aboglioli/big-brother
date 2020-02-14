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
	Update(id string, req *UpdateRequest) (*models.Organization, error)
	Delete(id string) error
}

// Request DTOs
type CreateRequest struct {
	Name string `json:"name" validate:"required,min=3,max=64"`
}

type UpdateRequest struct {
	Name string `json:"name" validate:"required,min=3,max=64"`
}

// Implementation
type service struct {
	repo      Repository
	validator Validator
}

func NewService(repo Repository) *service {
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
	if err := s.validator.CreateRequest(req); err != nil {
		return nil, err
	}

	org := models.NewOrganization()
	org.Name = req.Name

	if err := s.validator.Schema(org); err != nil {
		return nil, err
	}

	if err := s.repo.Insert(org); err != nil {
		return nil, errors.ErrInternalServer.Wrap(err)
	}

	return org, nil
}

func (s *service) Update(id string, req *UpdateRequest) (*models.Organization, error) {
	if err := s.validator.UpdateRequest(req); err != nil {
		return nil, err
	}

	org, err := s.getByID(id)
	if err != nil {
		return nil, err
	}

	org.Name = req.Name

	if err := s.validator.Schema(org); err != nil {
		return nil, err
	}

	if err := s.repo.Update(org); err != nil {
		return nil, errors.ErrInternalServer.Wrap(err)
	}

	return org, nil
}

func (s *service) Delete(id string) error {
	if _, err := s.getByID(id); err != nil {
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		return errors.ErrInternalServer.Wrap(err)
	}

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
