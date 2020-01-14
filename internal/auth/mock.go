package auth

import (
	"github.com/stretchr/testify/mock"
)

// Repository
type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) FindByID(tokenID string) (*Token, error) {
	args := m.Called(tokenID)
	return args.Get(0).(*Token), args.Error(1)
}

func (m *mockRepository) Insert(token *Token) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *mockRepository) Delete(tokenID string) error {
	args := m.Called(tokenID)
	return args.Error(0)
}

// Service
type mockService struct {
	*serviceImpl
	repo *mockRepository
}

func newMockService() *mockService {
	repo := &mockRepository{}
	serv := &serviceImpl{
		repo: repo,
	}
	return &mockService{serv, repo}
}
