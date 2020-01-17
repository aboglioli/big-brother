package auth

import (
	"github.com/aboglioli/big-brother/pkg/models"
	"github.com/stretchr/testify/mock"
)

// Encoder
type mockEncoder struct {
	mock.Mock
}

func (m *mockEncoder) Encode(tokenID string) (string, error) {
	args := m.Called(tokenID)
	return args.String(0), args.Error(1)
}

func (m *mockEncoder) Decode(tokenStr string) (string, error) {
	args := m.Called(tokenStr)
	return args.String(0), args.Error(1)
}

// Repository
type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) FindByID(tokenID string) (*models.Token, error) {
	args := m.Called(tokenID)
	if token, ok := args.Get(0).(*models.Token); ok {
		return token, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRepository) Insert(token *models.Token) error {
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
	enc  *mockEncoder
}

func newMockService() *mockService {
	repo := &mockRepository{}
	enc := &mockEncoder{}
	serv := &serviceImpl{
		repo: repo,
		enc:  enc,
	}
	return &mockService{serv, repo, enc}
}
