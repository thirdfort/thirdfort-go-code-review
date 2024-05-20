package repositories

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/thirdfort/thirdfort-go-code-review/internal"
	"github.com/thirdfort/thirdfort-go-code-review/internal/models"
)

// Create new Transaction with ID from PA
func (s *Store) CreateTransaction(ctx context.Context, transaction *models.Transaction) (*models.Transaction, error) {
	if transaction.CreatedAt.IsZero() {
		transaction.CreatedAt = time.Now()
	}
	transaction.UpdatedAt = time.Now()

	res := s.db.WithContext(ctx).Create(transaction)
	if res.Error != nil {
		s.logError(ctx, res)
		return nil, res.Error
	}

	txs, err := s.GetTransactions(ctx, &models.Transaction{ID: transaction.ID})
	if err != nil {
		return nil, err
	}

	return &txs[0], nil
}

func (s *Store) GetTransaction(ctx context.Context, transaction *models.Transaction) (*models.Transaction, error) {
	txs, err := s.GetTransactions(ctx, transaction)
	if err != nil {
		return nil, errors.Wrap(err, "GetTransaction")
	}

	if len(txs) > 0 {
		return &txs[0], nil
	}

	return nil, internal.ErrNotFound
}

func (s *Store) GetTransactions(ctx context.Context, transaction *models.Transaction) ([]models.Transaction, error) {
	var transactions []models.Transaction

	q := s.db.WithContext(ctx).Where(transaction)

	res := q.Find(&transactions)
	if res.Error != nil {
		s.logError(ctx, res)
		return nil, res.Error
	}

	return transactions, nil
}

func (s *Store) UpdateTransaction(ctx context.Context, transaction *models.Transaction) (*models.Transaction, error) {
	tx, err := s.GetTransaction(ctx, &models.Transaction{ID: transaction.ID})
	if err != nil {
		return nil, errors.Wrap(err, "GetTransaction")
	}

	if tx == nil {
		return nil, internal.ErrNotFound
	}

	res := s.db.WithContext(ctx).Save(transaction)
	if res.Error != nil {
		s.logError(ctx, res)
		return nil, res.Error
	}

	return s.GetTransaction(ctx, &models.Transaction{ID: transaction.ID})
}
