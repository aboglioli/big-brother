package organizations

import (
	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/models"
)

type Validator interface {
	Status(org *models.Organization) error
}

type validator struct {
}

func NewValidator() Validator {
	return &validator{}
}

func (v *validator) Status(org *models.Organization) error {
	if org == nil {
		return errors.ErrNotFound
	}
	if !org.Enabled {
		return errors.ErrNotFound
	}
	return nil
}
