package users

import (
	"github.com/aboglioli/big-brother/internal/auth"
	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/events"
)

// Errors
var (
	ErrNotFound     = errors.Status.New("user.service.not_found").S(404)
	ErrNotValidated = errors.Status.New("user.service.not_validated")
	ErrRegister     = errors.Status.New("user.service.register")
	ErrNotAvailable = errors.Validation.New("user.not_available")
	ErrUpdate       = errors.Status.New("user.service.update")
	ErrDelete       = errors.Status.New("user.service.delete")
	ErrInvalidUser  = errors.Status.New("user.service.invalid_user")
	ErrInvalidLogin = errors.Validation.New("user.service.invalid_login")
)

// Interfaces
type Service interface {
	GetByID(id string) (*User, error)

	Register(req *RegisterRequest) (*User, error)
	Update(id string, req *UpdateRequest) (*User, error)
	Delete(id string) error

	Login(req *LoginRequest) (*auth.Token, error)
	Logout(tokenStr string) error
}

// Implementations
type serviceImpl struct {
	repo      Repository
	events    events.Manager
	validator Validator
	authServ  auth.Service
}

func NewService(repo Repository, events events.Manager, authServ auth.Service) Service {
	return &serviceImpl{
		repo:      repo,
		events:    events,
		validator: NewValidator(),
		authServ:  authServ,
	}
}

func (s *serviceImpl) GetByID(id string) (*User, error) {
	return s.getByID(id)
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
}

func (s *serviceImpl) Register(req *RegisterRequest) (*User, error) {
	errs := make(errors.Errors, 0)

	// Password strength
	if err := s.validator.ValidatePassword(req.Password); err != nil {
		errs = append(errs, err)
	}

	user := NewUser()
	user.Username = req.Username
	user.SetPassword(req.Password)
	user.Email = req.Email
	user.Name = req.Name
	user.Lastname = req.Lastname

	// Schema validation
	if err := s.validator.ValidateSchema(user); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return nil, errs
	}

	// Is it available?
	vErr := ErrNotAvailable
	if existing, _ := s.repo.FindByUsername(req.Username); existing != nil {
		vErr = vErr.F("username", "not_available")
	}
	if existing, _ := s.repo.FindByEmail(req.Email); existing != nil {
		vErr = vErr.F("email", "not_available")
	}
	if len(vErr.Fields) > 0 {
		return nil, vErr
	}

	// Insert
	if err := s.repo.Insert(user); err != nil {
		return nil, ErrRegister.Wrap(err)
	}

	// Emit event
	userCreatedEvent := NewUserEvent(user, "UserCreated")
	if err := s.events.Publish(userCreatedEvent, &events.Options{"user", "user.created", ""}); err != nil {
		return nil, ErrRegister.Wrap(err)
	}

	return user, nil
}

type UpdateRequest struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
	Email    *string `json:"email"`
	Name     *string `json:"name"`
	Lastname *string `json:"lastname"`
}

func (s *serviceImpl) Update(id string, req *UpdateRequest) (*User, error) {
	user, err := s.getByID(id)
	if err != nil {
		return nil, err
	}

	errs := make(errors.Errors, 0)
	if req.Password != nil {
		if err := s.validator.ValidatePassword(*req.Password); err != nil {
			errs = append(errs, ErrPasswordValidation)
		} else {
			user.SetPassword(*req.Password)
		}
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Lastname != nil {
		user.Lastname = *req.Lastname
	}

	// Schema validation
	if err := s.validator.ValidateSchema(user); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return nil, errs
	}

	vErr := ErrNotAvailable
	if req.Username != nil {
		if existing, _ := s.repo.FindByUsername(*req.Username); existing != nil && existing.ID.Hex() != id {
			vErr = vErr.F("username", "not_available")
		} else {
			user.Username = *req.Username
		}
	}
	if req.Email != nil {
		if existing, _ := s.repo.FindByEmail(*req.Email); existing != nil && existing.ID.Hex() != id {
			vErr = vErr.F("email", "not_available")
		} else {
			user.Email = *req.Email
			user.Validated = false
		}
	}

	if len(vErr.Fields) > 0 {
		return nil, vErr
	}

	// Update
	if err := s.repo.Update(user); err != nil {
		return nil, ErrUpdate.C("id", id).Wrap(err)
	}

	// Emit event
	userUpdatedEvent := NewUserEvent(user, "UserUpdated")
	if err := s.events.Publish(userUpdatedEvent, &events.Options{"user", "user.updated", ""}); err != nil {
		return nil, ErrUpdate.Wrap(err)
	}

	return user, nil
}

func (s *serviceImpl) Delete(id string) error {
	user, err := s.getByID(id)
	if err != nil {
		return err
	}

	// Delete
	if err := s.repo.Delete(id); err != nil {
		return ErrDelete.Wrap(err)
	}

	// Emit event
	userDeletedEvent := NewUserEvent(user, "UserDeleted")
	if err := s.events.Publish(userDeletedEvent, &events.Options{"user", "user.deleted", ""}); err != nil {
		return ErrDelete.Wrap(err)
	}

	return nil
}

type LoginRequest struct {
	UsernameOrEmail *string `json:"username_or_email"`
	Password        *string `json:"password"`
}

func (s *serviceImpl) Login(req *LoginRequest) (*auth.Token, error) {
	vErr := ErrInvalidLogin
	if req.UsernameOrEmail == nil {
		vErr = vErr.F("username", "required")
	}
	if req.Password == nil {
		vErr = vErr.F("password", "required")
	}
	if len(vErr.Fields) > 0 {
		return nil, vErr
	}

	user, err := s.repo.FindByUsername(*req.UsernameOrEmail)
	if user == nil || err != nil {
		user, err = s.repo.FindByEmail(*req.UsernameOrEmail)
	}

	if user == nil || err != nil {
		return nil, ErrInvalidUser.Wrap(err)
	}

	if !user.ComparePassword(*req.Password) {
		return nil, ErrInvalidUser
	}

	token, err := s.authServ.Create(user.ID.Hex())
	if err != nil {
		return nil, ErrInvalidUser.Wrap(err)
	}

	return token, nil
}

func (s *serviceImpl) Logout(tokenStr string) error {
	token, err := s.authServ.Invalidate(tokenStr)

	if err != nil {
		return ErrInvalidUser.Wrap(err)
	}

	if existing, err := s.repo.FindByID(token.UserID); existing == nil || err != nil {
		return ErrInvalidUser.Wrap(err)
	}

	return nil
}

func (s *serviceImpl) getByID(id string) (*User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil || !user.Enabled {
		return nil, ErrNotFound.C("id", id).Wrap(err)
	}
	if !user.Validated {
		return nil, ErrNotValidated.C("id", id)
	}

	return user, nil
}
