package organizations

import (
	"testing"

	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/models"
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
