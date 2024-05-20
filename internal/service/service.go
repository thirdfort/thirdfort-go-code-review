package service

import (
	"context"

	"github.com/thirdfort/go-slogctx"
	"github.com/thirdfort/thirdfort-go-code-review/internal/cache"
	"github.com/thirdfort/thirdfort-go-code-review/internal/config"
	"github.com/thirdfort/thirdfort-go-code-review/internal/models"
	"github.com/thirdfort/thirdfort-go-code-review/internal/repositories"
)

type MainService interface {
	CheckTxOwnership(ctx context.Context, transactionID string) error
	GetActor(ctx context.Context) (*models.Actor, error)
	GetTransaction(ctx context.Context, tx *models.Transaction) (*models.Transaction, error)
	GetTransactions(ctx context.Context) ([]models.Transaction, error)
	PatchTransaction(ctx context.Context, txStatus *models.TransactionStatus) (*models.Transaction, error)
	Validate(item any, typeName string) (any, error)
}

type MockService struct {
	Logger    *slogctx.Logger
	DataStore repositories.DataStore
}

type Service struct {
	Logger    *slogctx.Logger
	DataStore repositories.DataStore
	cache     cache.Cache
}

func New(conf *config.Config,
	logger *slogctx.Logger,
	ds repositories.DataStore,
	cache cache.Cache,
) (*Service, error) {
	return &Service{
		Logger:    logger,
		DataStore: ds,
		cache:     cache,
	}, nil
}
