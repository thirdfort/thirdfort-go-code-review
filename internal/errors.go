package internal

import (
	"errors"
)

// nolint
var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")
	ErrBadRequest    = errors.New("bad request")
	ErrServerError   = errors.New("server error")
	ErrTooManyItems  = errors.New("too many items") // When expecting one but got more
	ErrUnsupported   = errors.New("unsupported")
)

// nolint
// database related errors
var (
	ErrDatabaseConnection = errors.New("database connection")
	ErrDatabaseMigration  = errors.New("database migration")
)
