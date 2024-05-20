package repositories

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/pkg/errors"
	"github.com/thirdfort/go-slogctx"
	"github.com/thirdfort/thirdfort-go-code-review/internal"
	"github.com/thirdfort/thirdfort-go-code-review/internal/models"
)

func (s *Store) CreateExpectationMetadata(ctx context.Context, exptMetadata *models.ExpectationMetadata) (*models.ExpectationMetadata, error) {
	res := s.db.WithContext(ctx).Create(exptMetadata.NewID())
	if res.Error != nil {
		s.logError(ctx, res)
		return nil, res.Error
	}

	met, err := s.GetExpectationMetadata(ctx, &models.ExpectationMetadata{ID: exptMetadata.ID})
	if err != nil {
		return nil, err
	}

	return met, nil
}

func (s *Store) GetExpectationMetadata(ctx context.Context, item *models.ExpectationMetadata) (*models.ExpectationMetadata, error) {
	var items []models.ExpectationMetadata

	q := s.db.WithContext(ctx).Where(item)

	res := q.Find(&items)
	if len(items) == 0 {
		return nil, internal.ErrNotFound
	}
	if len(items) != 1 {
		errStr := fmt.Sprintf("Found %d items of type %s, expected 1", len(items), item.TableName())
		slogctx.Error(ctx, errStr, slog.Any("item", item))
	}
	if res.Error != nil {
		s.logError(ctx, res)
		return nil, errors.Wrap(res.Error, "get")
	}

	if len(items) > 0 {
		retItem := items[0]
		return &retItem, nil
	}

	return nil, internal.ErrTooManyItems
}

func (s *Store) UpdateExpectationMetadata(ctx context.Context, exptMetadata *models.ExpectationMetadata) (*models.ExpectationMetadata, error) {
	res := s.db.WithContext(ctx).Updates(exptMetadata)
	if res.Error != nil {
		s.logError(ctx, res)
		return nil, res.Error
	}

	met, err := s.GetExpectationMetadata(ctx, &models.ExpectationMetadata{ID: exptMetadata.ID})
	if err != nil {
		return nil, err
	}

	return met, nil
}

func (s *Store) DeleteExpectationMetadata(ctx context.Context, exptMetadata *models.ExpectationMetadata) error {
	res := s.db.WithContext(ctx).Delete(exptMetadata)
	if res.Error != nil {
		s.logError(ctx, res)
		return res.Error
	}

	return nil
}
