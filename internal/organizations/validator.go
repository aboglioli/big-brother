package organizations

import (
	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/models"
	gValidator "github.com/aboglioli/big-brother/pkg/validator"
)

var (
	ErrSchemaValidation = errors.Validation.New("organization.invalid_schema")
)

type Validator interface {
	Status(org *models.Organization) error
	Schema(org *models.Organization) error
	CreateRequest(req *CreateRequest) error
	UpdateRequest(req *UpdateRequest) error
}

type validator struct {
	validator *gValidator.Validator
}

func NewValidator() *validator {
	return &validator{
		validator: gValidator.NewValidator(),
	}
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

func (v *validator) Schema(org *models.Organization) error {
	if err := v.validator.CheckFields(org); err != nil {
		return ErrSchemaValidation.Wrap(err)
	}
	return nil
}

func (v *validator) CreateRequest(req *CreateRequest) error {
	if err := v.validator.CheckFields(req); err != nil {
		return errors.ErrRequest.Wrap(err)
	}
	return nil
}

func (v *validator) UpdateRequest(req *UpdateRequest) error {
	if err := v.validator.CheckFields(req); err != nil {
		return errors.ErrRequest.Wrap(err)
	}
	return nil
}
