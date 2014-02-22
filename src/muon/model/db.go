package model

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"io"
)

var db *sql.DB

var NotFound = errors.New("model not found")

func Initialize(file string) (io.Closer, error) {
	var err error

	if db, err = sql.Open("sqlite3", file); err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	if err = createUsersTable(); err != nil {
		return nil, err
	}

	return db, nil
}
