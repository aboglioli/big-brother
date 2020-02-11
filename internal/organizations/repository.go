package organizations

import (
	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/models"
)

var (
	ErrRepositoryNotFound = errors.Internal.New("organization.repository.not_found")
	ErrRepositoryInsert   = errors.Internal.New("organization.repository.insert")
	ErrRepositoryUpdate   = errors.Internal.New("organization.repository.update")
	ErrRepositoryDelete   = errors.Internal.New("organization.repository.delete")
)

type Repository interface {
	FindByID(id string) (*models.Organization, error)
	SearchByName(name string) ([]*models.Organization, error)

	Insert(*models.Organization) error
	Update(*models.Organization) error
	Delete(id string) error
}
