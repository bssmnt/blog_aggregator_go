package database

import (
	"database/sql"
	"github.com/pressly/goose/v3"
	"os"
	"path/filepath"
)

func InitDB(dbURL string) (*Queries, error) {
	goose.SetLogger(goose.NopLogger())

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return nil, err
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	projectRoot := filepath.Join(currentDir, "..")

	migrationsPath := filepath.Join(projectRoot, "sql", "schema")

	if err := goose.Up(db, migrationsPath); err != nil {
		return nil, err
	}

	return New(db), nil
}
