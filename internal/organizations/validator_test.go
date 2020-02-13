package organizations

import (
	"testing"

	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/models"
	gValidator "github.com/aboglioli/big-brother/pkg/validator"
)

func TestValidatorStatus(t *testing.T) {
	mOrg := models.NewOrganization()

	type test struct {
		name string
		org  func() *models.Organization
		err  error
	}

	tests := []test{{
		"nil",
		func() *models.Organization {
			return nil
		},
		errors.ErrNotFound,
	}, {
		"not enabled",
		func() *models.Organization {
			org := mOrg.Clone()
			org.Enabled = false
			return org
		},
		errors.ErrNotFound,
	}, {
		"found",
		func() *models.Organization {
			org := mOrg.Clone()
			org.Enabled = true
			return org
		},
		nil,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			org := test.org()
			v := NewValidator()

			err := v.Status(org)
			errors.Assert(t, test.err, err)
		})
	}
}

func TestValidatorSchema(t *testing.T) {
	type test struct {
		name string
		org  *models.Organization
		err  error
	}

	tests := []test{{
		"empty organization",
		&models.Organization{
			Name: "",
		},
		ErrSchemaValidation,
	}, {
		"short name",
		&models.Organization{
			Name: "or",
		},
		ErrSchemaValidation,
	}, {
		"valid",
		&models.Organization{
			Name: "Org",
		},
		nil,
	}, {
		"valid with numbers",
		&models.Organization{
			Name: "Org 1",
		},
		nil,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := NewValidator()

			err := v.Schema(test.org)

			errors.Assert(t, test.err, err)
		})
	}
}

func TestValidatorCreateRequest(t *testing.T) {
	type test struct {
		name string
		req  *CreateRequest
		err  error
	}

	tests := []test{{
		"empty",
		&CreateRequest{
			Name: "",
		},
		errors.ErrRequest,
	}, {
		"short name: min 3",
		&CreateRequest{
			Name: "or",
		},
		errors.ErrRequest,
	}, {
		"valid",
		&CreateRequest{
			Name: "Organization 1",
		},
		nil,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := NewValidator()

			err := v.CreateRequest(test.req)

			errors.Assert(t, test.err, err)
		})
	}
}

func TestValidatorUpdateRequest(t *testing.T) {
	type test struct {
		name string
		req  *UpdateRequest
		err  error
	}

	tests := []test{{
		"empty",
		&UpdateRequest{
			Name: "",
		},
		errors.ErrRequest,
	}, {
		"short name",
		&UpdateRequest{
			Name: "or",
		},
		errors.ErrRequest.Wrap(gValidator.ErrFieldsValidation),
	}, {
		"valid",
		&UpdateRequest{
			Name: "Org 1",
		},
		nil,
	}}

	for _, test := range tests {
		v := NewValidator()

		err := v.UpdateRequest(test.req)

		errors.Assert(t, test.err, err)
	}
}
