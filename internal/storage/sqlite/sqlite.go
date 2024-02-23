package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3" // init sqlite3 driver
)

type Storage struct {
	db *sql.DB
}

func New(StoragePath string) (*Storage, error) {
	const fn = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", StoragePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	statement, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
	    id INTEGER PRIMARY KEY,
	    alias TEXT NOT NULL UNIQUE,
	    url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	_, err = statement.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &Storage{db: db}, nil
}
