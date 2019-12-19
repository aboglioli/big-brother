package users

import (
	"github.com/aboglioli/big-brother/errors"
	"github.com/aboglioli/big-brother/events"
)

// Errors
var (
	ErrNotFound     = errors.Status.New("user.service.not_found").S(404)
	ErrNotValidated = errors.Status.New("user.service.not_validated")
	ErrNotActive    = errors.Status.New("user.service.not_active")
	ErrCreate       = errors.Status.New("user.service.create")
	ErrNotAvailable = errors.Validation.New("user.not_available")
	ErrUpdate       = errors.Status.New("user.service.update")
	ErrDelete       = errors.Status.New("user.service.delete")
)

// Interfaces
type Service interface {
	GetByID(id string) (*User, error)

	Create(req *CreateRequest) (*User, error)
	Update(id string, req *UpdateRequest) (*User, error)
	Delete(id string) error
}

// Implementations
type serviceImpl struct {
	repo      Repository
	events    events.Manager
	validator Validator
}

func NewService(repo Repository, events events.Manager) Service {
	return &serviceImpl{
		repo:      repo,
		events:    events,
		validator: NewValidator(),
	}
}

func (s *serviceImpl) GetByID(id string) (*User, error) {
	return s.getByID(id)
}

type CreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (s *serviceImpl) Create(req *CreateRequest) (*User, error) {
	errs := make(errors.Errors, 0)

	// Is it available?
	vErr := ErrNotAvailable
	if existing, _ := s.repo.FindByUsername(req.Username); existing != nil {
		vErr = vErr.F("username", "not_available")
	}
	if existing, _ := s.repo.FindByEmail(req.Email); existing != nil {
		vErr = vErr.F("email", "not_available")
	}
	if len(vErr.Fields) > 0 {
		errs = append(errs, vErr)
	}

	// Password strength
	if err := s.validator.ValidatePassword(req.Password); err != nil {
		errs = append(errs, err)
	}

	// Create
	user := NewUser()
	user.Username = req.Username
	user.SetPassword(req.Password)
	user.Email = req.Email

	// Schema validation
	if err := s.validator.ValidateSchema(user); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return nil, errs
	}

	// Insert
	if err := s.repo.Insert(user); err != nil {
		return nil, ErrCreate.Wrap(err)
	}

	// Emit event
	userCreatedEvent := NewUserEvent(user)
	if err := s.events.Publish(userCreatedEvent, &events.Options{"user", "user.created", ""}); err != nil {
		return nil, ErrCreate.Wrap(err)
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
		errs = append(errs, vErr)
	}

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

	// Update
	if err := s.repo.Update(user); err != nil {
		return nil, ErrUpdate.C("id", id).Wrap(err)
	}

	// Emit event
	userUpdatedEvent := NewUserEvent(user)
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
	userDeletedEvent := NewUserEvent(user)
	if err := s.events.Publish(userDeletedEvent, &events.Options{"user", "user.deleted", ""}); err != nil {
		return ErrDelete.Wrap(err)
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
	if !user.Active {
		return nil, ErrNotActive.C("id", id)
	}

	return user, nil
}
