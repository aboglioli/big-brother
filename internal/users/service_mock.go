package users

import (
	"github.com/aboglioli/big-brother/mocks"
	"github.com/aboglioli/big-brother/pkg/models"
	"github.com/stretchr/testify/mock"
)

// Mocks
type mockValidator struct {
	mock.Mock
}

func (m *mockValidator) ValidateSchema(u *models.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *mockValidator) ValidatePassword(pwd string) error {
	args := m.Called(pwd)
	return args.Error(0)
}

type mockRepository struct {
	mock.Mock
}

func (r *mockRepository) FindByID(id string) (*models.User, error) {
	args := r.Called(id)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (r *mockRepository) FindByUsername(username string) (*models.User, error) {
	args := r.Called(username)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (r *mockRepository) FindByEmail(email string) (*models.User, error) {
	args := r.Called(email)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (r *mockRepository) Insert(u *models.User) error {
	args := r.Called(u)
	return args.Error(0)
}

func (r *mockRepository) Update(u *models.User) error {
	args := r.Called(u)
	return args.Error(0)
}

func (r *mockRepository) Delete(id string) error {
	args := r.Called(id)
	return args.Error(0)
}

// Auth service
type mockAuthService struct {
	mock.Mock
}

func (s *mockAuthService) Create(userID string) (string, error) {
	args := s.Called(userID)
	return args.String(0), args.Error(1)
}

func (s *mockAuthService) Validate(tokenStr string) (*models.Token, error) {
	args := s.Called(tokenStr)
	if token, ok := args.Get(0).(*models.Token); ok {
		return token, args.Error(1)
	}
	return nil, args.Error(1)
}

func (s *mockAuthService) Invalidate(tokenStr string) (*models.Token, error) {
	args := s.Called(tokenStr)
	if token, ok := args.Get(0).(*models.Token); ok {
		return token, args.Error(1)
	}
	return nil, args.Error(1)
}

// Service
type mockService struct {
	*serviceImpl
	repo      *mockRepository
	events    *mocks.MockEventManager
	validator *mockValidator
	authServ  *mockAuthService
}

func newMockService() *mockService {
	repo := &mockRepository{}
	events := mocks.NewMockEventManager()
	validator := &mockValidator{}
	authServ := &mockAuthService{}

	serv := &serviceImpl{
		repo:      repo,
		events:    events,
		validator: validator,
		authServ:  authServ,
	}

	return &mockService{serv, repo, events, validator, authServ}
}
