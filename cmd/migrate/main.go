package main

import (
	"database/sql"
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

var (
	glob = filepath.Join(migratePath, scriptsPattern)
)

type Migration struct {
	Database string
	Scripts  []string
}

func GetMigrations() []Migration {
	files, err := filepath.Glob(glob)
	if err != nil {
		panic(err)
	}

	m := make(map[string][]string)
	for _, file := range files {
		parts := strings.Split(file, "/")
		dbName := parts[1]
		m[dbName] = append(m[dbName], file)
	}

	var migrations []Migration
	for dbName, scripts := range m {
		migration := Migration{
			Database: dbName,
			Scripts:  scripts,
		}
		migrations = append(migrations, migration)
	}

	return migrations
}

func (m Migration) Run() {
	c := config.Get()

	db, err := db.ConnectPostgres(c.PostgresURL, m.Database, c.PostgresUsername, c.PostgresPassword)
	if err != nil {
		panic(err)
	}

	fmt.Printf("# Migrating %s:\n", m.Database)

	oldScripts := m.oldScripts(db)

	for _, script := range m.Scripts {
		if oldScripts[script] {
			fmt.Printf("- %s\n", script)
			continue
		}

		fmt.Printf("+ %s ", script)
		b, err := ioutil.ReadFile(script)
		if err != nil {
			panic(err)
		}

		sql := string(b)

		_, err = db.Exec(sql)
		if err != nil {
			fmt.Printf("ERROR\n\t%s\n", err)
			continue
		}

		m.addScript(db, script)

		fmt.Printf("OK\n")
	}

	fmt.Println()

	db.Close()
}

func (m Migration) addScript(db *sql.DB, script string) {
	_, err := db.Exec(`INSERT INTO migrations(script) VALUES($1)`, script)
	if err != nil {
		panic(err)
	}
}

func (m Migration) oldScripts(db *sql.DB) map[string]bool {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			script TEXT NOT NULL
		)
	`)
	if err != nil {
		panic(err)
	}

	rows, err := db.Query("SELECT script FROM migrations")
	if err != nil {
		panic(err)
	}

	scripts := make(map[string]bool)
	for rows.Next() {
		var s string
		err := rows.Scan(&s)
		if err != nil {
			panic(err)
		}
		scripts[s] = true
	}

	return scripts
}

func main() {
	migrations := GetMigrations()

	tm := Migration{
		Database: "test",
		Scripts:  make([]string, 0),
	}
	for _, m := range migrations {
		m.Run()
		tm.Scripts = append(tm.Scripts, m.Scripts...)
	}
	tm.Run()
}
