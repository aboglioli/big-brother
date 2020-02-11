package organizations

import (
	"testing"

	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestServiceGetByID(t *testing.T) {
	mOrg := models.NewOrganization()

	type test struct {
		name string
		id   string
		org  *models.Organization
		err  error
		mock func(m *mockService)
	}

	tests := []test{{
		"invalid id",
		"123",
		nil,
		errors.ErrNotFound,
		func(m *mockService) {
			m.repo.On("FindByID", "123").Return(nil, ErrRepositoryNotFound)
		},
	}, {
		"empty id",
		"",
		nil,
		errors.ErrNotFound,
		nil,
	}, {
		"valid id, non-existing organization",
		"abc-123",
		nil,
		errors.ErrNotFound,
		func(m *mockService) {
			m.repo.On("FindByID", "abc-123").Return(nil, ErrRepositoryNotFound)
		},
	}, {
		"not enabled organization",
		mOrg.ID,
		nil,
		errors.ErrNotFound,
		func(m *mockService) {
			c := mOrg.Clone()
			m.repo.On("FindByID", mOrg.ID).Return(c, nil)
			m.validator.On("Status", c).Return(errors.ErrNotFound)
		},
	}, {
		"existing organization",
		mOrg.ID,
		mOrg,
		nil,
		func(m *mockService) {
			c := mOrg.Clone()
			m.repo.On("FindByID", mOrg.ID).Return(c, nil)
			m.validator.On("Status", c).Return(nil)
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			serv := newMockService()
			if test.mock != nil {
				test.mock(serv)
			}
			org, err := serv.GetByID(test.id)

			if test.err != nil {
				errors.Assert(t, test.err, err)
				assert.Nil(org)
			} else {
				assert.Nil(err)
				if assert.NotNil(org) {
					assert.Equal(test.org, mOrg)
				}
			}
			serv.repo.AssertExpectations(t)
			serv.validator.AssertExpectations(t)
		})

	}
}

func TestServiceSearch(t *testing.T) {
	mOrg1 := models.NewOrganization()
	mOrg1.Name = "Organization 1"
	mOrg2 := models.NewOrganization()
	mOrg2.Name = "Organization 2"

	type test struct {
		name  string
		query string
		err   error
		orgs  []*models.Organization
		mock  func(m *mockService)
	}

	tests := []test{{
		"empty query",
		"",
		nil,
		[]*models.Organization{},
		nil,
	}, {
		"short query",
		"Org",
		nil,
		[]*models.Organization{mOrg1, mOrg2},
		func(m *mockService) {
			cOrg1, cOrg2 := mOrg1.Clone(), mOrg2.Clone()
			m.repo.On("SearchByName", "Org").Return([]*models.Organization{cOrg1, cOrg2}, nil)
		},
	}, {
		"long query",
		"Organization",
		nil,
		[]*models.Organization{mOrg1, mOrg2},
		func(m *mockService) {
			cOrg1, cOrg2 := mOrg1.Clone(), mOrg2.Clone()
			m.repo.On("SearchByName", "Organization").Return([]*models.Organization{cOrg1, cOrg2}, nil)
		},
	}, {
		"specific query",
		"Organization 2",
		nil,
		[]*models.Organization{mOrg2},
		func(m *mockService) {
			m.repo.On("SearchByName", "Organization 2").Return([]*models.Organization{mOrg2.Clone()}, nil)
		},
	}, {
		"error on searching",
		"Organization",
		errors.ErrNotFound,
		nil,
		func(m *mockService) {
			m.repo.On("SearchByName", "Organization").Return(nil, ErrRepositoryNotFound)
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			serv := newMockService()
			if test.mock != nil {
				test.mock(serv)
			}

			orgs, err := serv.Search(test.query)

			if test.err != nil {
				errors.Assert(t, test.err, err)
				assert.Nil(orgs)
			} else {
				assert.Nil(err)
				assert.Equal(test.orgs, orgs)
			}
		})
	}
}
