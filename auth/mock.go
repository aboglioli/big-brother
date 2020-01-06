package auth

import (
	"github.com/aboglioli/big-brother/mock"
)

// Repository
type mockRepository struct {
	Mock mock.Mock
	Repo *repositoryImpl
}

func newMockRepository() *mockRepository {
	cache := mock.NewMockCache("auth")
	return &mockRepository{
		Repo: &repositoryImpl{cache},
	}
}

func (m *mockRepository) FindByID(tokenID string) (*Token, error) {
	call := mock.Call("FindByID", tokenID)

	token, err := m.Repo.FindByID(tokenID)

	m.Mock.Called(call.Return(token, err))
	return token, err
}

func (m *mockRepository) Insert(token *Token) error {
	call := mock.Call("Insert", token)

	err := m.Repo.Insert(token)

	m.Mock.Called(call.Return(err))
	return err
}

func (m *mockRepository) Delete(tokenID string) error {
	call := mock.Call("Delete", tokenID)

	err := m.Repo.Delete(tokenID)

	m.Mock.Called(call.Return(err))
	return err
}

func (m *mockRepository) populate(tokens ...*Token) {
	for _, token := range tokens {
		m.Repo.Insert(token)
	}
}

// Service
type mockService struct {
	*serviceImpl
	repo *mockRepository
}

func newMockService() *mockService {
	repo := newMockRepository()
	serv := &serviceImpl{
		repo: repo,
	}
	return &mockService{serv, repo}
}
