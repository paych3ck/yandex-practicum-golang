package db

import (
	"database/sql"
	"errors"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL,
    title VARCHAR(255) NOT NULL,
    comment TEXT NOT NULL,
    repeat VARCHAR(128) NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

var conn *sql.DB

func Init(dbFile string) (*sql.DB, error) {
	install := false
	if _, err := os.Stat(dbFile); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			install = true
		} else {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	if install {
		if _, err := db.Exec(schema); err != nil {
			db.Close()
			return nil, err
		}
	}

	conn = db
	return conn, nil
}

func DB() *sql.DB {
	return conn
}
