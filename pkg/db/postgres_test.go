package db

import (
	"testing"

	"github.com/aboglioli/big-brother/pkg/config"
	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestPostgres(t *testing.T) {
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
			CREATE TABLE IF NOT EXISTS test_table (
				name text,
				age integer,
				lastname text
			);
		`)
		assert.Nil(err)

		_, err = db.Exec(`
			INSERT INTO test_table(name, age, lastname)
			VALUES('Alan', 21, NULL);
			INSERT INTO test_Table(name, age, lastname)
			VALUES('Boglioli', 23, 'Caffe');
		`)
		assert.Nil(err)

		rows, err := db.Query("SELECT COUNT(*) FROM test_table")
		assert.Nil(err)
		rows.Next()
		var count int
		err = rows.Scan(&count)
		assert.Nil(err)
		assert.Equal(2, count)

		type person struct {
			name     string
			age      int
			lastname *string
		}

		rows, err = db.Query("SELECT * FROM test_table")
		assert.Nil(err)
		people := make([]person, 0)
		for rows.Next() {
			var p person
			err := rows.Scan(&p.name, &p.age, &p.lastname)
			assert.Nil(err)
			people = append(people, p)
		}
		assert.Nil(rows.Err())
		assert.Len(people, 2)
		assert.Equal(person{"Alan", 21, nil}, people[0])
		assert.Equal(person{"Boglioli", 23, utils.NewString("Caffe")}, people[1])

		_, err = db.Exec("DROP TABLE test_table")
		assert.Nil(err)
	}

}
