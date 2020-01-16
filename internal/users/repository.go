package users

import (
	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/models"
)

// Errors
var (
	ErrRepositoryNotFound = errors.Internal.New("user.repository.not_found")
	ErrRepositoryInsert   = errors.Internal.New("user.repository.insert")
	ErrRepositoryUpdate   = errors.Internal.New("user.repository.update")
	ErrRepositoryDelete   = errors.Internal.New("user.repository.delete")
)

// Interfaces
type Repository interface {
	FindByID(id string) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)

	Insert(*models.User) error
	Update(*models.User) error
	Delete(id string) error
}
