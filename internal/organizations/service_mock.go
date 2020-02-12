package organizations

import (
	"github.com/aboglioli/big-brother/pkg/models"
	"github.com/stretchr/testify/mock"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) FindByID(id string) (*models.Organization, error) {
	args := m.Called(id)
	org, _ := args.Get(0).(*models.Organization)
	return org, args.Error(1)
}

func (m *mockRepository) SearchByName(name string) ([]*models.Organization, error) {
	args := m.Called(name)
	org, _ := args.Get(0).([]*models.Organization)
	return org, args.Error(1)
}

func (m *mockRepository) Insert(org *models.Organization) error {
	args := m.Called(org)
	return args.Error(0)
}

func (m *mockRepository) Update(org *models.Organization) error {
	args := m.Called(org)
	return args.Error(0)
}

func (m *mockRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

type mockValidator struct {
	mock.Mock
}

func (m *mockValidator) Status(org *models.Organization) error {
	args := m.Called(org)
	return args.Error(0)
}

func (m *mockValidator) Schema(org *models.Organization) error {
	args := m.Called(org)
	return args.Error(0)
}

func (m *mockValidator) CreateRequest(req *CreateRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *mockValidator) UpdateRequest(req *UpdateRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

type mockService struct {
	*service
	repo      *mockRepository
	validator *mockValidator
}

func newMockService() *mockService {
	repo := new(mockRepository)
	validator := new(mockValidator)
	serv := &service{
		repo:      repo,
		validator: validator,
	}
	return &mockService{serv, repo, validator}
}
