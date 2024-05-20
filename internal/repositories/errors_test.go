package repositories

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestIsRecordNotFoundTrue(t *testing.T) {
	err := gorm.ErrRecordNotFound

	assert.True(t, IsRecordNotFound(err))
}

func TestIsRecordNotFoundFalse(t *testing.T) {
	err := gorm.ErrDuplicatedKey

	assert.False(t, IsRecordNotFound(err))
}
