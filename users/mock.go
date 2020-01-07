package users

import (
	"time"

	"github.com/aboglioli/big-brother/auth"
	"github.com/aboglioli/big-brother/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User
const (
	userID  = "4af9f070eaf502a95c5271d4"
	adminID = "4af9f070eaf502a95c5271d5"
)

func newMockUser(id string) *User {
	user := NewUser()
	if id != "" {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			panic(err)
		}
		user.ID = objID
	}
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

// Auth service
type mockAuthService struct {
	Mock      mock.Mock
	tokensStr map[string]string
	tokens    map[string]*auth.Token
}

func newMockAuthService() *mockAuthService {
	return &mockAuthService{
		tokensStr: make(map[string]string),
		tokens:    make(map[string]*auth.Token),
	}
}

func (s *mockAuthService) Create(userID string) (*auth.Token, error) {
	call := mock.Call("Create", userID)

	tokenStr, ok1 := s.tokensStr[userID]
	token, ok2 := s.tokens[tokenStr]
	if !ok1 || !ok2 {
		s.Mock.Called(call.Return(nil, auth.ErrCreate))
		return nil, auth.ErrCreate
	}

	s.Mock.Called(call.Return(token, nil))
	return token, nil
}

func (s *mockAuthService) Validate(tokenStr string) (*auth.Token, error) {
	call := mock.Call("Validate", tokenStr)

	token, ok := s.tokens[tokenStr]
	if !ok {
		s.Mock.Called(call.Return(nil, auth.ErrUnauthorized))
		return nil, auth.ErrUnauthorized
	}

	s.Mock.Called(call.Return(token, nil))
	return token, nil

}

func (s *mockAuthService) Invalidate(tokenStr string) error {
	call := mock.Call("Invalidate", tokenStr)

	token, err := s.Validate(tokenStr)
	if err != nil {
		s.Mock.Called(call.Return(err))
		return err
	}

	delete(s.tokensStr, token.UserID)
	delete(s.tokens, tokenStr)

	s.Mock.Called(call.Return(nil))
	return nil
}

func (s *mockAuthService) populate(userIDs ...string) {
	for _, userID := range userIDs {
		token := auth.NewToken(userID)
		tokenStr, err := token.Encode()
		if err != nil {
			panic(err)
		}

		s.tokensStr[userID] = tokenStr
		s.tokens[tokenStr] = token
	}
}

// Service
type mockService struct {
	*serviceImpl
	repo      *mockRepository
	events    *mock.EventManager
	validator *mockValidator
	authServ  *mockAuthService
}

func newMockService() *mockService {
	repo := newMockRepository()
	events := mock.NewMockEventManager()
	validator := newMockValidator()
	authServ := newMockAuthService()
	serv := &serviceImpl{
		repo:      repo,
		events:    events,
		validator: validator,
		authServ:  authServ,
	}

	return &mockService{serv, repo, events, validator, authServ}
}
