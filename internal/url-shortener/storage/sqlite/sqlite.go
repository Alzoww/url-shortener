package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Alzoww/url-shortener/internal/url-shortener/storage"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const f = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", f, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS urls(
	    id INTEGER PRIMARY KEY,
	    url TEXT NOT NULL UNIQUE,
	    alias TEXT NOT NULL UNIQUE
	    );
	CREATE INDEX IF NOT EXISTS idx_alias ON urls(alias);
	    `)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", f, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", f, err)
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) SaveURL(urlToSave, alias string) error {
	const f = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare(`
	INSERT INTO urls(url, alias)
	VALUES ($1, $2)
`)
	if err != nil {
		return fmt.Errorf("%s: %w", f, err)
	}

	_, err = stmt.Exec(urlToSave, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s: %w", f, storage.ErrURLExists)
		}

		return fmt.Errorf("%s: %w", f, err)

	}

	return nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const f = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare(`
	SELECT url
	FROM urls
	WHERE alias = $1;
	`)
	if err != nil {
		return "", fmt.Errorf("%s: %w", f, err)
	}

	var resUrl string
	err = stmt.QueryRow(alias).Scan(&resUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", f, storage.ErrURLNotFound)
		}

		return "", fmt.Errorf("%s: %w", f, err)
	}

	return resUrl, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
