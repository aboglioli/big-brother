package db

import (
	"testing"

	"github.com/aboglioli/big-brother/pkg/config"
	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgres(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	db, err := ConnectPostgres("invalid", "invalid", "invalid", "invalid")
	assert.NotNil(t, err)
	assert.Nil(t, db)

	c := config.Get()
	db, err = ConnectPostgres(c.PostgresURL, "test", c.PostgresUsername, c.PostgresPassword)
	if !assert.Nil(t, err) {
		t.Log(err.(errors.Error).Cause)
	}
	require.NotNil(t, db)
	defer db.Close()

	_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS test_table (
				name text,
				age integer,
				lastname text
			);
		`)
	assert.Nil(t, err)

	_, err = db.Exec(`
		INSERT INTO test_table(name, age, lastname)
			VALUES('Alan', 21, NULL);
		INSERT INTO test_Table(name, age, lastname)
			VALUES('Boglioli', 23, 'Caffe');
	`)
	assert.Nil(t, err)

	type person struct {
		name     string
		age      int
		lastname *string
	}

	t.Run("count", func(t *testing.T) {
		assert := assert.New(t)
		row := db.QueryRow("SELECT COUNT(*) FROM test_table")
		assert.NotNil(row)
		var count int
		err = row.Scan(&count)
		assert.Nil(err)
		assert.Equal(2, count)
	})

	t.Run("multiple rows", func(t *testing.T) {
		assert := assert.New(t)
		rows, err := db.Query("SELECT * FROM test_table")
		defer rows.Close()
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
	})

	t.Run("query single row", func(t *testing.T) {
		assert := assert.New(t)
		row := db.QueryRow("SELECT * FROM test_table WHERE lastname = $1 and age > $2", "Caffe", 20)
		var p person
		err = row.Scan(&p.name, &p.age, &p.lastname)
		assert.Nil(err)
		assert.Equal(person{"Boglioli", 23, utils.NewString("Caffe")}, p)

		row = db.QueryRow("SELECT * FROM test_table WHERE lastname = $1 and age < $2", "Caffe", 20)
		err = row.Scan(&p.name, &p.age, &p.lastname)
		assert.NotNil(err)
	})

	t.Run("nil field", func(t *testing.T) {
		assert := assert.New(t)
		p := person{name: "name1", age: 25, lastname: utils.NewString("lastname1")}
		_, err = db.Exec("INSERT INTO test_table(name, age, lastname) VALUES($1, $2, $3)", &p.name, &p.age, &p.lastname)
		assert.Nil(err)

		p = person{}
		row := db.QueryRow("SELECT * FROM test_table WHERE name = 'name1'")
		err = row.Scan(&p.name, &p.age, &p.lastname)
		assert.Nil(err)
		assert.Equal(person{name: "name1", age: 25, lastname: utils.NewString("lastname1")}, p)

		p = person{name: "name2", age: 26, lastname: nil}
		_, err = db.Exec("INSERT INTO test_table(name, age, lastname) VALUES($1, $2, $3)", &p.name, &p.age, &p.lastname)
		assert.Nil(err)

		p = person{}
		row = db.QueryRow("SELECT * FROM test_table WHERE name = 'name2'")
		err = row.Scan(&p.name, &p.age, &p.lastname)
		assert.Nil(err)
		assert.Equal(person{name: "name2", age: 26, lastname: nil}, p)
	})

	t.Run("transaction", func(t *testing.T) {
		assert := assert.New(t)

		// Rollback
		tx, err := db.Begin()
		assert.Nil(err)
		assert.NotNil(tx)

		_, err = tx.Exec("INSERT INTO test_table(name, age, lastname) VALUES('Trans', 1, 'Action')")
		assert.Nil(err)

		var p person
		row := db.QueryRow("SELECT * FROM test_table WHERE name = 'Trans' AND lastname = 'Action'")
		err = row.Scan(&p.name, &p.age, &p.lastname)
		assert.NotNil(err)

		err = tx.Rollback()
		assert.Nil(err)

		row = db.QueryRow("SELECT * FROM test_table WHERE name = 'Trans' AND lastname = 'Action'")
		err = row.Scan(&p.name, &p.age, &p.lastname)
		assert.NotNil(err)

		// Commit
		tx, err = db.Begin()
		assert.Nil(err)
		assert.NotNil(tx)
		_, err = tx.Exec("INSERT INTO test_table(name, age, lastname) VALUES('Trans', 1, 'Action')")
		assert.Nil(err)
		err = tx.Commit()
		assert.Nil(err)

		row = db.QueryRow("SELECT * FROM test_table WHERE name = 'Trans' AND lastname = 'Action'")
		err = row.Scan(&p.name, &p.age, &p.lastname)
		assert.Nil(err)
		assert.Equal(person{"Trans", 1, utils.NewString("Action")}, p)
	})

	_, err = db.Exec("DROP TABLE test_table")
	assert.Nil(t, err)
}
