package repositories

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/thirdfort/thirdfort-go-code-review/internal"
	"github.com/thirdfort/thirdfort-go-code-review/internal/models"
	"gorm.io/gorm"
)

// Create new Actor with internal ID
func (s *Store) CreateFingerprint(ctx context.Context, fingerprint *models.Fingerprint) (*models.Fingerprint, error) {
	fingerprint.CreatedAt = time.Now()
	fingerprint.UpdatedAt = time.Now()

	res := s.db.WithContext(ctx).Create(fingerprint)
	if res.Error != nil {
		s.logError(ctx, res)
		return nil, res.Error
	}

	return s.GetFingerprint(ctx, &models.Fingerprint{Fingerprint: fingerprint.Fingerprint})
}

func (s *Store) GetFingerprint(ctx context.Context, fingerprint *models.Fingerprint) (*models.Fingerprint, error) {
	var ret models.Fingerprint

	q := s.db.WithContext(ctx).Model(fingerprint).Where(models.Fingerprint{Fingerprint: fingerprint.Fingerprint})

	res := q.First(&ret)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, internal.ErrNotFound
		}
		s.logError(ctx, res)
		return nil, res.Error
	}

	return &ret, nil
}
