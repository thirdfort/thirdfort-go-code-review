package repositories

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/pkg/errors"
	"github.com/thirdfort/go-slogctx"
	"github.com/thirdfort/thirdfort-go-code-review/internal"
	"github.com/thirdfort/thirdfort-go-code-review/internal/models"
	"gorm.io/gorm"
)

type Op[T models.DataType] struct {
	*Store
}

type MockOp[T models.DataType] interface {
	Create(ctx context.Context, item models.Item[T], filters map[string]any) (*T, error)
	Delete(ctx context.Context, item models.Item[T], filters map[string]any) error
	Get(ctx context.Context, item models.Item[T], filters map[string]any) (*T, error)
	Update(ctx context.Context, item models.Item[T], filters map[string]any) (*T, error)
}

func New[T models.DataType](s *Store) Op[T] {
	return Op[T]{s}
}

// Get a slice of items of type T
func (s Op[T]) Get(ctx context.Context, item models.Item[T], filters map[string]any) (*T, error) {
	var items []T
	q := s.db.WithContext(ctx).Where(item)

	for k, v := range filters {
		q.Where(fmt.Sprintf("%s = ?", k), v)
	}

	q.Order("created_at DESC")

	if item.TableName() == internal.DbTableDocument {
		q.Preload("Documents")
	}

	if item.TableName() == internal.DbTableExpectation && item.GetBase().Type == internal.TypeBankLink {
		q.Preload("Banks")
	}

	res := q.Debug().Find(&items)
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

// Insert a new row in database of type T
func (s Op[T]) Create(ctx context.Context, item models.Item[T], filters map[string]any) (*T, error) {
	var retItem T
	tx := s.db.WithContext(ctx).Begin()
	i := item.NewID()
	res := tx.Create(i).Scan(&retItem)
	if res.Error != nil {
		tx.Rollback()
		s.logError(ctx, res)
		return nil, res.Error
	}

	tx.Commit()

	return &retItem, nil
}

// Update row in database of type T
func (s Op[T]) Update(ctx context.Context, item models.Item[T], filters map[string]any) (*T, error) {
	q := s.db.WithContext(ctx)

	item.NewUpdated()

	res := q.Where(filters).Save(item)
	if res.Error != nil {
		s.logError(ctx, res)
		return nil, res.Error
	}

	return s.Get(ctx, item, nil)
}

// Update the changes in a row in database of type T
func (s Op[T]) UpdateChanges(ctx context.Context, item models.Item[T], filters map[string]any) (*T, error) {
	var retItem *T
	q := s.db.WithContext(ctx)

	item.NewUpdated()

	res := q.Where(filters).Updates(item).Scan(&retItem)
	if res.Error != nil {
		s.logError(ctx, res)
		return nil, res.Error
	}

	return retItem, nil
}

// Delete row in database of type T
func (s Op[T]) Delete(ctx context.Context, item models.Item[T], filters map[string]any) error {
	q := s.db.WithContext(ctx)

	for k, v := range filters {
		q.Where(fmt.Sprintf("%s = ?", k), v)
	}

	res := q.Delete(item)
	if res.Error != nil {
		s.logError(ctx, res)
		return res.Error
	}

	return nil
}

// Get item with just the IDs set
func (s Op[T]) getBaseItem(item models.Item[T]) models.Item[T] {
	var vType any
	oldBase := item.GetBase()
	base := models.Base{
		Type:          oldBase.Type,
		TransactionID: oldBase.TransactionID,
		Shared: models.Shared{ // we should always get not deleted task
			DeletedAt: gorm.DeletedAt{},
		},
	}

	vType = item.SetBase(base)

	return vType.(models.Item[T])
}
