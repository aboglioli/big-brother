package users

import (
	"github.com/aboglioli/big-brother/internal/auth"
	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/events"
	"github.com/aboglioli/big-brother/pkg/models"
)

// Errors
var (
	ErrUserNotAvailable = errors.Status.New("user.not_available")
	ErrInvalidUser      = errors.Status.New("user.invalid_user")
)

// Interfaces
type Service interface {
	GetByID(id string) (*models.User, error)

	Register(req *RegisterRequest) (*models.User, error)
	Update(id string, req *UpdateRequest) (*models.User, error)
	ChangePassword(id string, req *ChangePasswordRequest) error
	Delete(id string) error

	Login(req *LoginRequest) (string, error)
	Logout(tokenStr string) error
}

// Request DTOs
type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required"`
	Lastname string `json:"lastname" validate:"required"`
}

type UpdateRequest struct {
	Username *string `json:"username"`
	Email    *string `json:"email" validate:"email"`
	Name     *string `json:"name"`
	Lastname *string `json:"lastname"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required"`
}

type LoginRequest struct {
	UsernameOrEmail string `json:"username" validate:"required"`
	Password        string `json:"password" validate:"required"`
}

// Implementations
type service struct {
	repo      Repository
	events    events.Bus
	validator Validator
	crypt     PasswordCrypt
	authServ  auth.Service
}

func NewService(repo Repository, events events.Bus, authServ auth.Service) Service {
	return &service{
		repo:      repo,
		events:    events,
		validator: NewValidator(),
		crypt:     NewBcryptCrypt(),
		authServ:  authServ,
	}
}

func (s *service) GetByID(id string) (*models.User, error) {
	return s.getByID(id)
}

func (s *service) Register(req *RegisterRequest) (*models.User, error) {
	// Request
	if err := s.validator.RegisterRequest(req); err != nil {
		return nil, err
	}

	// Password strength
	if err := s.validator.Password(req.Password); err != nil {
		return nil, err
	}

	user := models.NewUser()
	user.Username = req.Username
	user.Password = req.Password
	user.Email = req.Email
	user.Name = req.Name
	user.Lastname = req.Lastname

	// Schema validation
	if err := s.validator.Schema(user); err != nil {
		return nil, err
	}

	// Is it available?
	aErr := ErrUserNotAvailable
	if existing, _ := s.repo.FindByUsername(req.Username); existing != nil {
		aErr = aErr.F("username", "not_available")
	}
	if existing, _ := s.repo.FindByEmail(req.Email); existing != nil {
		aErr = aErr.F("email", "not_available")
	}
	if len(aErr.Fields) > 0 {
		return nil, aErr
	}

	// Set password
	hash, err := s.crypt.Hash(req.Password)
	if err != nil {
		return nil, errors.ErrInternalServer.Wrap(err)
	}
	user.Password = hash

	// Insert
	if err := s.repo.Insert(user); err != nil {
		return nil, errors.ErrInternalServer.Wrap(err)
	}

	return user, nil
}

func (s *service) Update(id string, req *UpdateRequest) (*models.User, error) {
	if err := s.validator.UpdateRequest(req); err != nil {
		return nil, err
	}

	user, err := s.getByID(id)
	if err != nil {
		return nil, err
	}

	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Lastname != nil {
		user.Lastname = *req.Lastname
	}

	// Schema validation
	if err := s.validator.Schema(user); err != nil {
		return nil, err
	}

	// Available
	aErr := ErrUserNotAvailable
	if req.Username != nil {
		if existing, _ := s.repo.FindByUsername(*req.Username); existing != nil && existing.ID != id {
			aErr = aErr.F("username", "not_available")
		} else {
			user.Username = *req.Username
		}
	}
	if req.Email != nil {
		if existing, _ := s.repo.FindByEmail(*req.Email); existing != nil && existing.ID != id {
			aErr = aErr.F("email", "not_available")
		} else {
			user.Email = *req.Email
			user.Validated = false
		}
	}
	if len(aErr.Fields) > 0 {
		return nil, aErr
	}

	// Update
	if err := s.repo.Update(user); err != nil {
		return nil, errors.ErrInternalServer.Wrap(err)
	}

	return user, nil
}

func (s *service) ChangePassword(id string, req *ChangePasswordRequest) error {
	if err := s.validator.ChangePasswordRequest(req); err != nil {
		return err
	}

	user, err := s.getByID(id)
	if err != nil {
		return err
	}

	if err := s.validator.Password(req.NewPassword); err != nil {
		return err
	}

	if !s.crypt.Compare(user.Password, req.CurrentPassword) {
		return ErrInvalidUser
	}

	hash, err := s.crypt.Hash(req.NewPassword)
	if err != nil {
		return errors.ErrInternalServer.Wrap(err)
	}

	user.Password = hash
	if err := s.repo.Update(user); err != nil {
		return errors.ErrInternalServer.Wrap(err)
	}

	return nil
}

func (s *service) Delete(id string) error {
	_, err := s.getByID(id)
	if err != nil {
		return err
	}

	// Delete
	if err := s.repo.Delete(id); err != nil {
		return errors.ErrInternalServer.Wrap(err)
	}

	return nil
}

func (s *service) Login(req *LoginRequest) (string, error) {
	if err := s.validator.LoginRequest(req); err != nil {
		return "", err
	}

	user, err := s.repo.FindByUsername(req.UsernameOrEmail)
	if user == nil || err != nil {
		user, err = s.repo.FindByEmail(req.UsernameOrEmail)
	}

	if user == nil || err != nil {
		return "", ErrInvalidUser.Wrap(err)
	}

	if !s.crypt.Compare(user.Password, req.Password) {
		return "", ErrInvalidUser
	}

	tokenStr, err := s.authServ.Create(user.ID)
	if err != nil {
		return "", ErrInvalidUser.Wrap(err)
	}

	return tokenStr, nil
}

func (s *service) Logout(tokenStr string) error {
	token, err := s.authServ.Invalidate(tokenStr)
	if err != nil {
		return ErrInvalidUser.Wrap(err)
	}

	if existing, err := s.repo.FindByID(token.UserID); existing == nil || err != nil {
		return ErrInvalidUser.Wrap(err)
	}

	return nil
}

func (s *service) getByID(id string) (*models.User, error) {
	if id == "" {
		return nil, errors.ErrNotFound.M("invalid id")
	}

	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.ErrNotFound.Wrap(err)
	}
	if err := s.validator.Status(user); err != nil {
		return nil, err
	}

	return user, nil
}
