package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/thirdfort/thirdfort-go-code-review/internal/models"
)

func (s *Store) FindAddress(ctx context.Context, transactionId string) (*models.Address, error) {
	txn := models.Transaction{}
	s.db.Where(fmt.Sprintf("id = %v", transactionId)).First(&txn)

	if txn.Address == nil {
		return nil, errors.New("no transaction")
	} else {
		return txn.Address, nil
	}
}

func (s *Store) UpdateAddress(ctx context.Context, transactionID string, newAddress models.Address) (*models.Address, error) {
	db := s.db
	txn := models.Transaction{}
	db.Where(fmt.Sprintf("id = %v", transactionID)).First(&txn)

	if txn.Address == nil {
		return nil, errors.New("no transaction")
	} else {
		txn.Address = &newAddress
		return txn.Address, nil
	}
}
