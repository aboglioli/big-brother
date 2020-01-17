package db

import (
	"testing"

	"github.com/aboglioli/big-brother/pkg/config"
	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	assert := assert.New(t)

	db, err := ConnectPostgres("invalid", "invalid", "invalid", "invalid")
	assert.NotNil(err)
	assert.Nil(db)

	c := config.Get()
	db, err = ConnectPostgres(c.PostgresURL, "test", c.PostgresUsername, c.PostgresPassword)
	if !assert.Nil(err) {
		t.Log(err.(errors.Error).Cause)
	}
	if assert.NotNil(db) {
		_, err := db.Exec(`
			DROP TABLE test_table;
			CREATE TABLE test_table (
				name text,
				age integer
			);
		`)
		assert.Nil(err)

		_, err = db.Exec(`
			INSERT INTO test_table(name, age)
			VALUES('Alan', 21);
			INSERT INTO test_Table(name, age)
			VALUES('Boglioli', 23);
		`)
		assert.Nil(err)

		rows, err := db.Query("SELECT COUNT(*) FROM test_table")
		assert.Nil(err)
		for rows.Next() {
			var count int
			err := rows.Scan(&count)
			assert.Nil(err)
			assert.Equal(count, 2)
		}

		type person struct {
			name string
			age  int
		}

		rows, err = db.Query("SELECT * FROM test_table")
		assert.Nil(err)
		people := make([]person, 0)
		for rows.Next() {
			var p person
			err := rows.Scan(&p.name, &p.age)
			assert.Nil(err)
			people = append(people, p)
		}
		assert.Nil(rows.Err())
		assert.Len(people, 2)
	}

}
