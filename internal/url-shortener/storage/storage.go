package storage

import "errors"

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists   = errors.New("url exists")
)

type Interface interface {
	SaveURL(urlToSave, alias string) error
	GetURL(alias string) (string, error)
}
