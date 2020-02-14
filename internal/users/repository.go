package users

import (
	"database/sql"
	"fmt"
	"time"

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

const (
	userFields = "id, username, password, email, name, lastname, role, validated, enabled, created_at, updated_at, deleted_at"
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
	return r.findOne(fmt.Sprintf("SELECT %s FROM users WHERE id = $1 AND enabled = true", userFields), id)
}

func (r *sqlRepository) FindByUsername(username string) (*models.User, error) {
	return r.findOne(fmt.Sprintf("SELECT %s FROM users WHERE username = $1 AND enabled = true", userFields), username)
}

func (r *sqlRepository) FindByEmail(email string) (*models.User, error) {
	return r.findOne(fmt.Sprintf("SELECT %s FROM users WHERE email = $1 AND enabled = true", userFields), email)
}

func (r *sqlRepository) Insert(u *models.User) error {
	u.CreatedAt = time.Now().UTC()

	sql := fmt.Sprintf(`
		INSERT INTO users(%s)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`, userFields)
	res, err := r.db.Exec(
		sql,
		u.ID,
		u.Username,
		u.Password,
		u.Email,
		u.Name,
		u.Lastname,
		u.Role,
		u.Validated,
		u.Enabled,
		u.CreatedAt,
		u.UpdatedAt,
		u.DeletedAt,
	)
	if err != nil {
		return ErrRepositoryInsert.Wrap(err)
	}
	if c, err := res.RowsAffected(); c == 0 || err != nil {
		return ErrRepositoryInsert
	}

	return nil
}

func (r *sqlRepository) Update(u *models.User) error {
	u.UpdatedAt = &time.Time{}
	*u.UpdatedAt = time.Now().UTC()

	res, err := r.db.Exec(`
		UPDATE users
		SET username = $2,
		password = $3,
		email = $4,
		name = $5,
		lastname = $6,
		role = $7,
		validated = $8,
		updated_at = $9
		WHERE id = $1
	`, u.ID, u.Username, u.Password, u.Email, u.Name, u.Lastname, u.Role, u.Validated, u.UpdatedAt)
	if err != nil {
		return ErrRepositoryUpdate.Wrap(err)
	}
	if c, err := res.RowsAffected(); c == 0 || err != nil {
		return ErrRepositoryUpdate
	}
	return nil
}

func (r *sqlRepository) Delete(id string) error {
	res, err := r.db.Exec("UPDATE users SET enabled = false, deleted_at = $2 WHERE id = $1", id, time.Now().UTC())
	if err != nil {
		return ErrRepositoryDelete.Wrap(err)
	}
	if c, err := res.RowsAffected(); c == 0 || err != nil {
		return ErrRepositoryDelete
	}
	return nil
}

func (r *sqlRepository) findOne(sql string, args ...interface{}) (*models.User, error) {
	row := r.db.QueryRow(sql, args...)

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
