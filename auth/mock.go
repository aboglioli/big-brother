package auth

import (
	"github.com/aboglioli/big-brother/cache"
)

// Service
type mockService struct {
	*serviceImpl
	repo *repositoryImpl
}

func newMockService() *mockService {
	cache := cache.NewInMemory("auth")
	repo := &repositoryImpl{
		cache: cache,
	}
	serv := &serviceImpl{
		repo: repo,
	}
	return &mockService{serv, repo}
}
