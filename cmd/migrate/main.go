package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/aboglioli/big-brother/pkg/config"
	"github.com/aboglioli/big-brother/pkg/db"
)

const (
	migratePath    = "./migrations"
	scriptsPattern = "/*/*.sql"
)

func listFiles(dir string) ([]string, error) {
	files, err := filepath.Glob(fmt.Sprintf("%s/%s", dir, scriptsPattern))
	return files, err
}

func readFile(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	return string(b), err
}

func parseDBName(path string) string {
	parts := strings.Split(path, "/")
	return parts[1]
}

func main() {
	files, err := listFiles(migratePath)
	if err != nil {
		panic(err)
	}

	scripts := map[string][]string{}
	for _, file := range files {
		db := parseDBName(file)
		if db != "" {
			scripts[db] = append(scripts[db], file)
		}
	}

	c := config.Get()

	testDB, err := db.ConnectPostgres(c.PostgresURL, "test", c.PostgresUsername, c.PostgresPassword)
	defer func() {
		testDB.Close()
	}()
	if err != nil {
		panic(err)
	}

	for database, files := range scripts {
		fmt.Printf("Populating %s:\n", database)
		db, err := db.ConnectPostgres(c.PostgresURL, database, c.PostgresUsername, c.PostgresPassword)
		if err != nil {
			panic(err)
		}

		_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS migrations (
				script TEXT NOT NULL
			)
		`)
		if err != nil {
			panic(err)
		}

		for _, file := range files {
			fmt.Printf("+ %s ", file)
			sql, err := readFile(file)
			if err != nil {
				panic(err)
			}

			testDB.Exec(sql)

			_, err = db.Exec(sql)
			if err != nil {
				fmt.Printf("\n\tERROR: %s\n", err)
				continue
			}
			_, err = db.Exec(`
				INSERT INTO migrations(script)
					VALUES($1)
			`, file)
			if err != nil {
				panic(err)
			}

			fmt.Printf("OK\n")

		}

		db.Close()
	}
}
