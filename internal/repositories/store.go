package repositories

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/DATA-DOG/go-sqlmock"
	slogGorm "github.com/orandin/slog-gorm"
	"github.com/thirdfort/go-slogctx"
	"github.com/thirdfort/thirdfort-go-code-review/internal"
	"github.com/thirdfort/thirdfort-go-code-review/internal/config"
	"github.com/thirdfort/thirdfort-go-code-review/internal/migrate"
	"github.com/thirdfort/thirdfort-go-code-review/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DataStore interface {
	GetStore() *Store
	Rollback(tx *gorm.DB)
	CreateActor(ctx context.Context, actorData map[string]string) (*models.Actor, error)
	GetActor(ctx context.Context, actor *models.Actor) (*models.Actor, error)
	UpdateActor(ctx context.Context, actorID *string, actorData map[string]string) (*models.Actor, error)
	CreateTransaction(ctx context.Context, transaction *models.Transaction) (*models.Transaction, error)
	UpdateTransaction(ctx context.Context, transaction *models.Transaction) (*models.Transaction, error)
	GetTransaction(ctx context.Context, transaction *models.Transaction) (*models.Transaction, error)
	GetTransactions(ctx context.Context, transaction *models.Transaction) ([]models.Transaction, error)
}

// We'll inject this dependency into the handlers and feed it down to the repositories where it'll be used
type Store struct {
	conf *config.Config
	db   *gorm.DB
	log  *slogctx.Logger
}

type MockStore struct {
	db  *gorm.DB
	log *slogctx.Logger
}

func NewMockStore(log *slogctx.Logger) (*Store, sqlmock.Sqlmock, error) {
	gormLogger := slogGorm.New(
		slogGorm.WithLogger(log.GetLogger()),
		slogGorm.SetLogLevel(slogGorm.DefaultLogType, slog.LevelDebug),
		slogGorm.WithTraceAll(),
	)

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Error(nil, fmt.Sprintf("Error : %s when opening a stub database connection", err), slog.Any("err", err))
		return nil, nil, err
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		Logger:         gormLogger,
		TranslateError: true,
		QueryFields:    false,
	})
	if err != nil {
		log.Error(nil, fmt.Sprintf("Error : %s when opening gorm database", err))
		return nil, nil, err
	}

	return &Store{
		db:  gormDB,
		log: log,
	}, mock, nil
}

func NewStore(conf *config.Config, log *slogctx.Logger) (*Store, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=UTC",
		conf.Postgres.Host,
		conf.Postgres.User,
		conf.Postgres.Password,
		conf.Postgres.Database,
		conf.Postgres.Port)

	return NewStoreFromDSN(conf, dsn, log)
}

func NewStoreFromDSN(conf *config.Config, dsn string, log *slogctx.Logger) (*Store, error) {
	pgConfig := postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}

	opts := []slogGorm.Option{
		slogGorm.WithLogger(log.GetLogger()),
		slogGorm.SetLogLevel(slogGorm.DefaultLogType, slog.LevelDebug),
		slogGorm.WithTraceAll(),
	}

	if conf.Postgres.EnableDebug {
		opts = append(opts, slogGorm.WithTraceAll())
	}

	gormLogger := slogGorm.New(
		opts...,
	)

	gormConfig := gorm.Config{
		Logger:         gormLogger,
		TranslateError: true,
		QueryFields:    true,
	}

	db, err := gorm.Open(postgres.New(pgConfig), &gormConfig)
	if err != nil {
		log.Error(nil, "error connecting to db", slog.Any("err", err))
		return nil, internal.ErrDatabaseConnection
	}

	if err := migrate.Migrate(db); err != nil {
		log.Error(nil, "migration error", slog.Any("err", err))
		return nil, internal.ErrDatabaseMigration
	}

	slogctx.Info(nil, "Connected to database",
		slog.String("host", conf.Postgres.Host),
		slog.String("port", conf.Postgres.Port),
		slog.String("database", conf.Postgres.Database),
	)

	return &Store{
		db:   db,
		log:  log,
		conf: conf,
	}, nil
}

func (s *Store) Close() {
	db, err := s.db.DB()
	if err == nil {
		db.Close()
	}
}

func (s *Store) CreateDatabase(dbName string) *gorm.DB {
	if s.conf.App.Env == config.Testing || s.conf.App.Env == config.Local {
		return s.db.Exec(fmt.Sprintf("CREATE DATABASE \"%s\"", dbName))
	}

	slogctx.Error(nil, "Not allowed to create database outside testing environment",
		slog.String("environment", string(s.conf.App.Env)))

	return nil
}

func (s *Store) DropTables() error {
	if s.conf.App.Env == config.Testing || s.conf.App.Env == config.Local {
		return migrate.DropTables(s.db)
	}

	err := errors.New("Not allowed to drop tables outside testing environment")

	slogctx.Error(nil, err.Error(),
		slog.String("environment", string(s.conf.App.Env)))

	return err
}

// TODO: Check how to get sql - this doesn't seem to work
func (s *Store) logError(ctx context.Context, res *gorm.DB) {
	slogctx.Error(ctx, "error performing query",
		slog.String("sql", res.Statement.SQL.String()),
		slog.Any("vars", res.Statement.Vars),
		slog.Any("err", res.Error.Error()))
}

func (s *Store) GetStore() *Store {
	return s
}

func (s *Store) Rollback(tx *gorm.DB) {
	if r := recover(); r != nil {
		tx.Rollback()
		panic(r)
	} else if err := recover(); err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
}
