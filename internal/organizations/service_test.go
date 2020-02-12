package organizations

import (
	"testing"

	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
				assert.Equal(test.org, mOrg)
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
			serv.repo.AssertExpectations(t)
		})
	}
}

func TestServiceCreate(t *testing.T) {
	type test struct {
		name string
		req  *CreateRequest
		err  error
		org  *models.Organization
		mock func(m *mockService)
	}

	tests := []test{{
		"empty request",
		&CreateRequest{
			Name: "",
		},
		errors.ErrRequest,
		nil,
		func(m *mockService) {
			m.validator.On("CreateRequest", &CreateRequest{
				Name: "",
			}).Return(errors.ErrRequest)
		},
	}, {
		"valid request, wrong schema",
		&CreateRequest{
			Name: "Org_1",
		},
		ErrSchemaValidation,
		nil,
		func(m *mockService) {
			m.validator.On("CreateRequest", &CreateRequest{
				Name: "Org_1",
			}).Return(nil)
			m.validator.On("Schema", mock.MatchedBy(func(org *models.Organization) bool {
				return org.Name == "Org_1"
			})).Return(ErrSchemaValidation)
		},
	}, {
		"valid request",
		&CreateRequest{
			Name: "Organization 1",
		},
		nil,
		&models.Organization{
			Name: "Organization 1",
		},
		func(m *mockService) {
			m.validator.On("CreateRequest", &CreateRequest{
				Name: "Organization 1",
			}).Return(nil)
			m.validator.On("Schema", mock.MatchedBy(func(org *models.Organization) bool {
				return org.Name == "Organization 1"
			})).Return(nil)
			m.repo.On("Insert", mock.MatchedBy(func(org *models.Organization) bool {
				return org.Name == "Organization 1" && org.ID != ""
			})).Return(nil)
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			serv := newMockService()
			if test.mock != nil {
				test.mock(serv)
			}

			org, err := serv.Create(test.req)

			if test.err != nil {
				errors.Assert(t, test.err, err)
				assert.Nil(org)
			} else {
				assert.Nil(err)
				if assert.NotNil(org) {
					assert.Equal(test.org.Name, org.Name)
					assert.NotEmpty(org.ID)
				}
			}
			serv.repo.AssertExpectations(t)
			serv.validator.AssertExpectations(t)
		})
	}
}

func TestServiceUpdate(t *testing.T) {
	mOrg := models.NewOrganization()
	mOrg.Name = "Organization 1"

	type test struct {
		name string
		id   string
		req  *UpdateRequest
		err  error
		org  *models.Organization
		mock func(m *mockService)
	}

	tests := []test{{
		"empty",
		"123",
		&UpdateRequest{
			Name: "",
		},
		errors.ErrRequest,
		nil,
		func(m *mockService) {
			m.validator.On("UpdateRequest", &UpdateRequest{
				Name: "",
			}).Return(errors.ErrRequest)
		},
	}, {
		"invalid id",
		"123",
		&UpdateRequest{
			Name: "New Organization 1",
		},
		errors.ErrNotFound,
		nil,
		func(m *mockService) {
			m.validator.On("UpdateRequest", &UpdateRequest{
				Name: "New Organization 1",
			}).Return(nil)
			m.repo.On("FindByID", "123").Return(nil, ErrRepositoryNotFound)
		},
	}, {
		"not enabled",
		"123",
		&UpdateRequest{
			Name: "New Organization 1",
		},
		errors.ErrNotFound,
		nil,
		func(m *mockService) {
			m.validator.On("UpdateRequest", &UpdateRequest{
				Name: "New Organization 1",
			}).Return(nil)
			m.repo.On("FindByID", "123").Return(mOrg.Clone(), nil)
			m.validator.On("Status", mOrg).Return(errors.ErrNotFound)
		},
	}, {
		"valid update",
		"123",
		&UpdateRequest{
			Name: "New Organization 1",
		},
		nil,
		&models.Organization{
			Name: "New Organization 1",
		},
		func(m *mockService) {
			m.validator.On("UpdateRequest", &UpdateRequest{
				Name: "New Organization 1",
			}).Return(nil)
			m.repo.On("FindByID", "123").Return(mOrg.Clone(), nil)
			m.validator.On("Status", mOrg).Return(nil)
			m.validator.On("Schema", mock.MatchedBy(func(org *models.Organization) bool {
				return org.Name == "New Organization 1"
			})).Return(nil)
			m.repo.On("Update", mock.MatchedBy(func(org *models.Organization) bool {
				return org.Name == "New Organization 1"
			})).Return(nil)
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			s := newMockService()
			if test.mock != nil {
				test.mock(s)
			}

			org, err := s.Update(test.id, test.req)

			if test.err != nil {
				errors.Assert(t, test.err, err)
				assert.Nil(org)
			} else {
				assert.Nil(err)
				if assert.NotNil(org) {
					assert.Equal(test.req.Name, org.Name)
				}
			}
			s.repo.AssertExpectations(t)
			s.validator.AssertExpectations(t)
		})
	}
}

func TestServiceDelete(t *testing.T) {
	mOrg := models.NewOrganization()
	mOrg.Name = "Organization 1"

	type test struct {
		name string
		id   string
		err  error
		mock func(m *mockService)
	}

	tests := []test{{
		"empty id",
		"",
		errors.ErrNotFound,
		nil,
	}, {
		"non existing org",
		"123",
		errors.ErrNotFound,
		func(m *mockService) {
			m.repo.On("FindByID", "123").Return(nil, ErrRepositoryNotFound)
		},
	}, {
		"non enabled org",
		"123",
		errors.ErrNotFound,
		func(m *mockService) {
			c := mOrg.Clone()
			m.repo.On("FindByID", "123").Return(c, nil)
			m.validator.On("Status", c).Return(errors.ErrNotFound)
		},
	}, {
		"existing org",
		"123",
		nil,
		func(m *mockService) {
			c := mOrg.Clone()
			m.repo.On("FindByID", "123").Return(c, nil)
			m.validator.On("Status", c).Return(nil)
			m.repo.On("Delete", "123").Return(nil)
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := newMockService()
			if test.mock != nil {
				test.mock(s)
			}

			err := s.Delete(test.id)

			errors.Assert(t, test.err, err)
		})
	}
}
