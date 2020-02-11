package employees

import (
	"context"

	"github.com/aboglioli/big-brother/pkg/models"
)

type Service interface {
	GetByID(emplID string) (*models.Employee, error)
	GetByUser(userID string) ([]*models.Employee, error)
	GetByOrganization(orgID string) ([]*models.Employee, error)
	GetByUserAndOrganization(userID, orgID string) (*models.Employee, error)

	Create(ctx context.Context, req *CreateRequest) (*models.Employee, error)
	Update(ctx context.Context, req *UpdateRequest) (*models.Employee, error)
	Delete(ctx context.Context, id string) (*models.Employee, error)
}

type CreateRequest struct {
	Name   string `json:"name"`
	RoleID string `json:"role_id"`
}

type UpdateRequest struct {
	Name string `json:"name"`
}
