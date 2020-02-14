package users

import (
	"database/sql"
	"testing"
	"time"

	"github.com/aboglioli/big-brother/pkg/config"
	"github.com/aboglioli/big-brother/pkg/db"
	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/models"
	"github.com/stretchr/testify/assert"
)

func populate() (*sql.DB, []*models.User) {
	c := config.Get()
	db, err := db.ConnectPostgres(c.PostgresURL, "test", c.PostgresUsername, c.PostgresPassword)
	db.Exec("DELETE FROM users")
	if err != nil {
		panic(err)
	}

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
	users := []*models.User{user1, user2, user3}

	for _, user := range users {
		_, err := db.Exec(`
			INSERT INTO users(id, username, password, email, name, lastname, role, validated, enabled, created_at)
				VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`, user.ID, user.Username, user.Password, user.Email, user.Name, user.Lastname, user.Role, user.Validated, user.Enabled, user.CreatedAt)
		if err != nil {
			panic(err)
		}
	}

	return db, users
}

func TestRepositoryFindByID(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	db, users := populate()
	defer db.Close()

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
		"not enabled",
		users[2].ID,
		ErrRepositoryNotFound,
		nil,
	}, {
		"existing",
		users[0].ID,
		nil,
		users[0],
	}, {
		"another",
		users[1].ID,
		nil,
		users[1],
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

func TestRepositoryFindByUsername(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	db, users := populate()
	defer db.Close()

	tests := []struct {
		name     string
		username string
		err      error
		user     *models.User
	}{{
		"empty username",
		"",
		ErrRepositoryNotFound,
		nil,
	}, {
		"short username",
		"u",
		ErrRepositoryNotFound,
		nil,
	}, {
		"not enabled",
		users[2].Username,
		ErrRepositoryNotFound,
		nil,
	}, {
		"existing user",
		users[0].Username,
		nil,
		users[0],
	}, {
		"existing user",
		users[1].Username,
		nil,
		users[1],
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			r := NewRepository(db)

			user, err := r.FindByUsername(test.username)

			if test.err != nil {
				errors.Assert(t, test.err, err)
				assert.Nil(user)
			} else {
				assert.Nil(err)
				if assert.NotNil(user) {
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
			}
		})
	}
}

func TestRepositoryFindByEmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	db, users := populate()
	defer db.Close()

	tests := []struct {
		name  string
		email string
		err   error
		user  *models.User
	}{{
		"empty email",
		"",
		ErrRepositoryNotFound,
		nil,
	}, {
		"invalid email",
		"q.com",
		ErrRepositoryNotFound,
		nil,
	}, {
		"short email",
		"u@u.com",
		ErrRepositoryNotFound,
		nil,
	}, {
		"not enabled",
		users[2].Email,
		ErrRepositoryNotFound,
		nil,
	}, {
		"existing user",
		users[0].Email,
		nil,
		users[0],
	}, {
		"existing user",
		users[1].Email,
		nil,
		users[1],
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			r := NewRepository(db)

			user, err := r.FindByEmail(test.email)

			if test.err != nil {
				errors.Assert(t, test.err, err)
				assert.Nil(user)
			} else {
				assert.Nil(err)
				if assert.NotNil(user) {
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
			}
		})
	}
}

func TestRepositoryInsert(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	db, users := populate()
	defer db.Close()

	user4 := models.NewUser()
	user4.Username = "user4"
	user4.Password = "hashed.password.4"
	user4.Email = "user4@user.com"
	user4.Name = "Fourth"
	user4.Lastname = "User"

	tests := []struct {
		name string
		user *models.User
		err  error
	}{{
		"empty user",
		models.NewUser(),
		nil,
	}, {
		"existing user",
		users[0],
		ErrRepositoryInsert,
	}, {
		"valid user",
		user4,
		nil,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			r := NewRepository(db)

			err := r.Insert(test.user)

			if test.err != nil {
				errors.Assert(t, test.err, err)
			} else {
				row := db.QueryRow("SELECT id, username, password, email, created_at FROM users WHERE id = $1", test.user.ID)
				var id, username, password, email string
				var createdAt time.Time
				row.Scan(&id, &username, &password, &email, &createdAt)
				assert.Equal(test.user.ID, id)
				assert.Equal(test.user.Username, username)
				assert.Equal(test.user.Password, password)
				assert.Equal(test.user.Email, email)
				assert.Equal(test.user.CreatedAt.Format("2006-01-02T15:04:05-0700"), createdAt.Format("2006-01-02T15:04:05-0700"))
			}
		})
	}
}

func TestRepositoryUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	db, users := populate()
	defer db.Close()

	user4 := models.NewUser()
	user4.Username = "user4"
	user4.Password = "hashed.password.4"
	user4.Email = "user4@user.com"
	user4.Name = "Fourth"
	user4.Lastname = "User"

	users[0].Username = "new-username"
	users[0].Name = "New Name"
	users[0].Lastname = "New Lastname"
	users[0].Validated = true

	tests := []struct {
		name string
		user *models.User
		err  error
	}{{
		"non existing user",
		user4,
		ErrRepositoryUpdate,
	}, {
		"existing user",
		users[0],
		nil,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			r := NewRepository(db)

			err := r.Update(test.user)

			if test.err != nil {
				errors.Assert(t, test.err, err)
			} else {
				assert.Nil(err)
				row := db.QueryRow("SELECT id, username, name, lastname, validated, updated_at FROM users WHERE id = $1", test.user.ID)
				var id, username, name, lastname string
				var validated bool
				var updatedAt time.Time
				row.Scan(&id, &username, &name, &lastname, &validated, &updatedAt)
				assert.Equal(test.user.ID, id)
				assert.Equal(test.user.Username, username)
				assert.Equal(test.user.Name, name)
				assert.Equal(test.user.Lastname, lastname)
				assert.Equal(test.user.Validated, validated)
				assert.NotEmpty(updatedAt)
			}
		})
	}
}

func TestRepositoryDelete(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	db, users := populate()

	tests := []struct {
		name string
		id   string
		err  error
	}{{
		"invalid id",
		"user123",
		ErrRepositoryDelete,
	}, {
		"non existing user",
		models.NewID(),
		ErrRepositoryDelete,
	}, {
		"existing user",
		users[0].ID,
		nil,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			r := NewRepository(db)

			err := r.Delete(test.id)

			if test.err != nil {
				errors.Assert(t, test.err, err)
			} else {
				row := db.QueryRow("SELECT enabled, deleted_at FROM users WHERE id = $1", test.id)
				var enabled bool
				var deletedAt time.Time
				row.Scan(&enabled, &deletedAt)
				assert.Equal(false, enabled)
				assert.NotEmpty(deletedAt)
			}
		})
	}
}
