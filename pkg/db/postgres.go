package db

import (
	"database/sql"
	"fmt"

	"github.com/aboglioli/big-brother/pkg/errors"
	_ "github.com/lib/pq"
)

var (
	ErrPostgresConnect = errors.Internal.New("postgres.connect")
)

func ConnectPostgres(url, database, username, password string) (*sql.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", username, password, url, database)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, ErrPostgresConnect.Wrap(err)
	}

	if err := db.Ping(); err != nil {
		return nil, ErrPostgresConnect.Wrap(err)
	}

	return db, nil
}
