package repositories

import (
	"errors"

	"gorm.io/gorm"
)

// Simple helper to determine whether error is a NotFound error, something explicit we'd like to report on
func IsRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
