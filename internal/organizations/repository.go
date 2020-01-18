package organizations

import "github.com/aboglioli/big-brother/pkg/models"

type Repository interface {
	FindByID(id string) (*models.Organization, error)
	Search(name string) ([]*models.Organization, error)

	Insert(*models.Organization) error
	Update(*models.Organization) error
	Delete(id string) error
}
