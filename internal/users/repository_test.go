package users

import (
	"database/sql"
	"testing"

	"github.com/aboglioli/big-brother/pkg/config"
	"github.com/aboglioli/big-brother/pkg/db"
	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func populate(db *sql.DB, users ...*models.User) error {
	for _, user := range users {
		_, err := db.Exec(`
			INSERT INTO users(id, username, password, email, name, lastname, role, validated, enabled, created_at)
				VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`, user.ID, user.Username, user.Password, user.Email, user.Name, user.Lastname, user.Role, user.Validated, user.Enabled, user.CreatedAt)
		if err != nil {
			return err
		}
	}

	return nil
}

func TestRepositoryFindByID(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	c := config.Get()
	db, err := db.ConnectPostgres(c.PostgresURL, "test", c.PostgresUsername, c.PostgresPassword)
	require.Nil(t, err)
	require.NotNil(t, db)
	defer func() {
		db.Exec("DELETE FROM users")
	}()

	user1 := models.NewUser()
	user1.Username = "user1"
	user1.Password = "hashed.password.1"
	user1.Email = "user1@user.com"
	user1.Name = "First"
	user1.Lastname = "User"
	user2 := models.NewUser()
	user2.Username = "user2"
	user2.Password = "hashed.password.2"
	user2.Email = "user2@user.com"
	user2.Name = "Second"
	user2.Lastname = "User"
	user2.Validated = true
	user3 := models.NewUser()
	user3.Username = "user3"
	user3.Password = "hashed.password.3"
	user3.Email = "user3@user.com"
	user3.Name = "Third"
	user3.Lastname = "User"
	user3.Validated = true
	user3.Enabled = false

	err = populate(db, user1, user2, user3)
	require.Nil(t, err)

	tests := []struct {
		name string
		id   string
		err  error
		user *models.User
	}{{
		"invalid id",
		"user123",
		ErrRepositoryNotFound,
		nil,
	}, {
		"non existing",
		"ed84840c-6e9c-46c4-9409-86283a9fa961",
		ErrRepositoryNotFound,
		nil,
	}, {
		"existing",
		user2.ID,
		nil,
		user2.Clone(),
	}, {
		"another",
		user3.ID,
		nil,
		user3.Clone(),
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			repo := NewRepository(db)
			user, err := repo.FindByID(test.id)

			if test.err != nil {
				errors.Assert(t, test.err, err)
				assert.Nil(user)
			} else {
				if !assert.Nil(err) {
					if err, ok := err.(errors.Error); ok {
						t.Errorf("%v", err.Cause)
					}
				}
				assert.Equal(test.user.ID, user.ID)
				assert.Equal(test.user.Username, user.Username)
				assert.Equal(test.user.Password, user.Password)
				assert.Equal(test.user.Email, user.Email)
				assert.Equal(test.user.Name, user.Name)
				assert.Equal(test.user.Lastname, user.Lastname)
				assert.Equal(test.user.Role, user.Role)
				assert.Equal(test.user.Validated, user.Validated)
				assert.Equal(test.user.Enabled, user.Enabled)
				assert.Equal(test.user.CreatedAt.Format("2006-01-02T15:04:05-0700"), user.CreatedAt.Format("2006-01-02T15:04:05-0700"))
			}
		})
	}
}
