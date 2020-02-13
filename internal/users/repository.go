package users

import (
	"database/sql"

	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/models"
)

// Errors
var (
	ErrRepositoryNotFound = errors.Internal.New("user.repository.not_found")
	ErrRepositoryInsert   = errors.Internal.New("user.repository.insert")
	ErrRepositoryUpdate   = errors.Internal.New("user.repository.update")
	ErrRepositoryDelete   = errors.Internal.New("user.repository.delete")
)

// Interfaces
type Repository interface {
	FindByID(id string) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)

	Insert(*models.User) error
	Update(*models.User) error
	Delete(id string) error
}

type sqlRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &sqlRepository{
		db: db,
	}
}

func (r *sqlRepository) FindByID(id string) (*models.User, error) {
	row := r.db.QueryRow(`
		SELECT id, username, password, email, name, lastname, role, validated, enabled, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1
	`, id)

	var u models.User
	err := row.Scan(
		&u.ID,
		&u.Username,
		&u.Password,
		&u.Email,
		&u.Name,
		&u.Lastname,
		&u.Role,
		&u.Validated,
		&u.Enabled,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeletedAt,
	)
	if err != nil {
		return nil, ErrRepositoryNotFound.Wrap(err)
	}

	return &u, nil
}

func (r *sqlRepository) FindByUsername(username string) (*models.User, error) {
	return nil, nil
}

func (r *sqlRepository) FindByEmail(email string) (*models.User, error) {
	return nil, nil
}

func (r *sqlRepository) Insert(u *models.User) error {
	return nil
}

func (r *sqlRepository) Update(u *models.User) error {
	return nil
}

func (r *sqlRepository) Delete(id string) error {
	return nil
}
