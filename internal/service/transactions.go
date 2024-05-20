package service

import (
	"context"
	"log/slog"

	"github.com/pkg/errors"
	"github.com/thirdfort/go-slogctx"
	"github.com/thirdfort/thirdfort-go-code-review/internal"
	"github.com/thirdfort/thirdfort-go-code-review/internal/models"
)

func (s *Service) PatchTransaction(ctx context.Context, txStatus *models.TransactionStatus) (*models.Transaction, error) {
	actor, err := s.GetActor(ctx)
	if err != nil {
		return nil, err
	}

	t := &models.Transaction{
		ID:      &txStatus.ID,
		ActorID: actor.ID,
	}

	// Get the original transaction from db for the eventID
	tx, err := s.GetTransaction(ctx, t)
	if err != nil {
		return nil, err
	}
	slogctx.Debug(nil, "db", slog.Any("tx", tx))

	transaction, err := s.DataStore.UpdateTransaction(ctx, tx)
	if err != nil {
		return nil, errors.Wrap(err, "PatchTransaction: UpdateTransaction")
	}

	if transaction == nil {
		return nil, internal.ErrNotFound
	}

	return tx, nil
}

func (s *Service) GetTransaction(ctx context.Context, tx *models.Transaction) (*models.Transaction, error) {
	actor, err := s.GetActor(ctx)
	if err != nil {
		return nil, internal.ErrForbidden
	}
	transaction, err := s.DataStore.GetTransaction(ctx, &models.Transaction{ID: tx.ID, ActorID: actor.ID})
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (s *Service) CheckTxOwnership(ctx context.Context, transactionID string) error {
	actor, err := s.GetActor(ctx)
	if err != nil || actor == nil || actor.ID == nil {
		return internal.ErrNotFound
	}

	transaction, err := s.DataStore.GetTransaction(ctx, &models.Transaction{ID: &transactionID, ActorID: actor.ID})
	if err != nil || transaction == nil || transaction.ID == nil {
		return internal.ErrNotFound
	}

	return nil
}

func (s *Service) GetTransactions(ctx context.Context) ([]models.Transaction, error) {
	// TODO: Add filtering by 'status'
	var transactions []models.Transaction
	actor, err := s.GetActor(ctx)
	if err != nil {
		// We might not have retrieved transactions+actor from PA so this error is allowed,
		// logging for debugging purposes
		slogctx.Warn(ctx, "Could not get actor", slog.Any("err", err))
	} else {
		transactions, err = s.DataStore.GetTransactions(ctx, &models.Transaction{ActorID: actor.ID})
		if err != nil {
			return nil, err
		}
	}

	return transactions, nil
}

func (s *Service) getTxFromSlice(txID *string, transactions []models.Transaction) *models.Transaction {
	for _, tx := range transactions {
		if *tx.ID == *txID {
			return &tx
		}
	}

	return nil
}

func updateTransactions(transactions []models.Transaction, updateTx *models.Transaction) []models.Transaction {
	for i, tx := range transactions {
		if tx.ID == updateTx.ID {
			transactions[i] = *updateTx
		}
	}

	return transactions
}
