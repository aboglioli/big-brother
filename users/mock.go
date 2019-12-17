package users

import (
	"time"

	"github.com/aboglioli/big-brother/tools/mock"
)

// Repository
type mockRepository struct {
	Mock       mock.Mock
	Collection []*User
}

func NewMockRepository() *mockRepository {
	return &mockRepository{
		Mock:       mock.Mock{},
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
	r.Collection = append(r.Collection, u)

	r.Mock.Called(call.Return(nil))
	return nil
}

func (r *mockRepository) Update(u *User) error {
	call := mock.Call("Update", u)

	for i, item := range r.Collection {
		if item.ID.Hex() == u.ID.Hex() {
			u.UpdatedAt = time.Now()
			r.Collection[i] = u
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
