package users

import (
	"github.com/aboglioli/big-brother/errors"
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
	FindByID(id string) (*User, error)
	FindByUsername(username string) (*User, error)
	FindByEmail(email string) (*User, error)

	Insert(*User) error
	Update(*User) error
	Delete(id string) error
}
