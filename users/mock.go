package users

import (
	"time"

	"github.com/aboglioli/big-brother/mock"
)

// User
func newMockUser() *User {
	user := NewUser()
	user.Username = "user"
	user.SetPassword("123456789")
	user.Email = "user@user.com"
	user.Name = "Name"
	user.Lastname = "Lastname"
	return user
}

// Validator
type mockValidator struct {
	Mock      mock.Mock
	validator Validator
}

func (m *mockValidator) ValidateSchema(u *User) error {
	call := mock.Call("ValidateSchema", u)
	err := m.validator.ValidateSchema(u)
	m.Mock.Called(call.Return(err))
	return err
}

func (m *mockValidator) ValidatePassword(pwd string) error {
	call := mock.Call("ValidatePassword", pwd)
	err := m.validator.ValidatePassword(pwd)
	m.Mock.Called(call.Return(err))
	return err
}

func newMockValidator() *mockValidator {
	return &mockValidator{
		validator: NewValidator(),
	}
}

// Repository
type mockRepository struct {
	Mock       mock.Mock
	Collection []*User
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		Collection: make([]*User, 0),
	}
}

func (r *mockRepository) FindByID(id string) (*User, error) {
	call := mock.Call("FindByID", id)

	for _, item := range r.Collection {
		if item.ID.Hex() == id {
			r.Mock.Called(call.Return(item, nil))
			return item, nil
		}
	}

	r.Mock.Called(call.Return(nil, ErrRepositoryNotFound))
	return nil, ErrRepositoryNotFound
}

func (r *mockRepository) FindByUsername(username string) (*User, error) {
	call := mock.Call("FindByUsername", username)

	for _, item := range r.Collection {
		if item.Username == username {
			r.Mock.Called(call.Return(item, nil))
			return item, nil
		}
	}

	r.Mock.Called(call.Return(nil, ErrRepositoryNotFound))
	return nil, ErrRepositoryNotFound
}

func (r *mockRepository) FindByEmail(email string) (*User, error) {
	call := mock.Call("FindByEmail", email)

	for _, item := range r.Collection {
		if item.Email == email {
			r.Mock.Called(call.Return(item, nil))
			return item, nil
		}
	}

	r.Mock.Called(call.Return(nil, ErrRepositoryNotFound))
	return nil, ErrRepositoryNotFound
}

func (r *mockRepository) Insert(u *User) error {
	call := mock.Call("Insert", u)

	for _, item := range r.Collection {
		if item.ID.Hex() == u.ID.Hex() {
			r.Mock.Called(call.Return(ErrRepositoryInsert))
			return ErrRepositoryInsert
		}
	}

	u.CreatedAt = time.Now()
	r.Collection = append(r.Collection, copyUser(u))

	r.Mock.Called(call.Return(nil))
	return nil
}

func (r *mockRepository) Update(u *User) error {
	call := mock.Call("Update", u)

	for i, item := range r.Collection {
		if item.ID.Hex() == u.ID.Hex() {
			u.UpdatedAt = time.Now()
			r.Collection[i] = copyUser(u)
			r.Mock.Called(call.Return(nil))
			return nil
		}
	}

	r.Mock.Called(call.Return(ErrRepositoryUpdate))
	return ErrRepositoryUpdate
}

func (r *mockRepository) Delete(id string) error {
	call := mock.Call("Delete", id)

	for _, item := range r.Collection {
		if item.ID.Hex() == id {
			item.DeletedAt = time.Now()
			item.Enabled = false
			return nil
		}
	}

	r.Mock.Called(call.Return(ErrRepositoryDelete))
	return ErrRepositoryDelete
}

func (r *mockRepository) populate(users ...*User) {
	r.Collection = make([]*User, 0)
	for _, user := range users {
		r.Collection = append(r.Collection, copyUser(user))
	}
}

func copyUser(u *User) *User {
	copy := *u
	return &copy
}

// Service
type mockService struct {
	*serviceImpl
	repo      *mockRepository
	events    *mock.EventManager
	validator *mockValidator
}

func newMockService() *mockService {
	repo := newMockRepository()
	events := mock.NewMockEventManager()
	validator := newMockValidator()
	serv := &serviceImpl{
		repo:      repo,
		events:    events,
		validator: validator,
	}

	return &mockService{serv, repo, events, validator}
}
