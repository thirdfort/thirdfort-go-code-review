package service

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/thirdfort/go-slogctx"
	"github.com/thirdfort/thirdfort-go-code-review/internal/cache"
	"github.com/thirdfort/thirdfort-go-code-review/internal/config"
	"github.com/thirdfort/thirdfort-go-code-review/internal/models"
	"github.com/thirdfort/thirdfort-go-code-review/internal/repositories"
)

type MainService interface {
	CheckTxOwnership(ctx context.Context, transactionID string) error
	GetActor(ctx context.Context) (*models.Actor, error)
	GetPersonalInformation(ctx context.Context, item *models.PersonalInformation) (*models.PersonalInformation, error)
	GetTasks(ctx context.Context, txID string) []any
	GetTransaction(ctx context.Context, tx *models.Transaction) (*models.Transaction, error)
	GetTransactions(ctx context.Context) ([]models.Transaction, error)
	HandlePutExpectation(ctx context.Context, transactionID *string, expectationID *string, item any) error
	PatchTransaction(ctx context.Context, txStatus *models.TransactionStatus) (*models.Transaction, error)
	PutPersonalInformation(ctx context.Context, item *models.PersonalInformation) (*models.PersonalInformation, error)
	Validate(item any, typeName string) (any, error)
}

type Ops[T models.DataType] interface {
	Create(context.Context, T) (T, error)
	Get(context.Context, T) (any, error)
	Update(context.Context, T) (T, error)
	Delete(context.Context, T) error
}

type MockService struct {
	Logger    *slogctx.Logger
	DataStore repositories.DataStore
	validator *validator.Validate
}

type Service struct {
	Logger    *slogctx.Logger
	DataStore repositories.DataStore
	validator *validator.Validate
	cache     cache.Cache
}

func New(conf *config.Config,
	logger *slogctx.Logger,
	ds repositories.DataStore,
	cache cache.Cache,
) (*Service, error) {
	dsn := conf.Sentry.Dsn
	if conf.App.Debug {
		dsn = ""
	}

	if conf.Sentry.Environment != config.Local {
		err := sentry.Init(sentry.ClientOptions{
			Debug:       conf.Sentry.Environment == config.Local,
			Dsn:         dsn,
			Environment: string(conf.Sentry.Environment),
		})
		if err != nil {
			return nil, errors.Wrap(err, "initialising sentry")
		}
	}

	return &Service{
		Logger:    logger,
		DataStore: ds,
		validator: validator.New(),
		cache:     cache,
	}, nil
}
