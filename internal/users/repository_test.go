package users

import (
	"database/sql"
	"testing"
	"time"

	"github.com/aboglioli/big-brother/pkg/models"
)

func populate(db *sql.DB) error {
	if _, err := db.Exec(`
		INSERT INTO users(id, username, password, email, name, lastname, role, validated, enabled, created_at)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, models.NewID(), "user1", "hashed.password.1", "user1@user.com", "User", "One", "user", true, true, time.Now()); err != nil {
		return err
	}

	return nil
}

func TestFindByID(t *testing.T) {
	// c := config.Get()
	// db, err := db.ConnectPostgres(c.PostgresURL, "users_and_organizations", c.PostgresUsername, c.PostgresPassword)
	// require.Nil(t, err)
	// require.NotNil(t, db)

	// err = populate(db)
	// require.Nil(t, err)
}
