package cache

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/thirdfort/go-slogctx"
	"github.com/thirdfort/thirdfort-go-code-review/internal"
	"github.com/thirdfort/thirdfort-go-code-review/internal/models"
)

type data struct {
	actorData *models.Actor
	txData    map[string]map[string]any
	ttl       time.Time
}

type Store struct {
	lock sync.Mutex
	data map[string]data
}

type Cache interface {
	GetActor(fingerprint string) *models.Actor
	SetActor(fingerprint string, actorData *models.Actor)
	GetTxData(fingerprint, transactionID, key string) any
	SetTxData(fingerprint, transactionID, key string, value any) error
	ResetTxData(fingerprint, transactionID, key string)
}

func New() *Store {
	return &Store{
		data: make(map[string]data),
	}
}

// Delete a key-value pair from the transaction data
func (a *Store) ResetTxData(fingerprint, transactionID, key string) {
	a.purge()
	a.lock.Lock()
	defer a.lock.Unlock()

	dataVal, ok := a.data[fingerprint]
	if !ok {
		return
	}

	txData, ok := dataVal.txData[transactionID]
	if !ok {
		return
	}

	_, ok = txData[key]
	if !ok {
		return
	}

	delete(txData, key)

	dataVal.txData[transactionID] = txData
	dataVal.ttl = time.Now().Add(internal.ActorCacheTTL)
	a.data[fingerprint] = dataVal
}

// Get a transaction related value by key
func (a *Store) GetTxData(fingerprint, transactionID, key string) any {
	a.purge()
	a.lock.Lock()
	defer a.lock.Unlock()

	dataVal, ok := a.data[fingerprint]
	if !ok {
		return nil
	}

	txData, ok := dataVal.txData[transactionID]
	if !ok {
		return nil
	}

	val, ok := txData[key]
	if !ok {
		return nil
	}

	return val
}

// Store a transaction related key-value -pair
func (a *Store) SetTxData(fingerprint, transactionID, key string, value any) error {
	a.purge()
	a.lock.Lock()
	defer a.lock.Unlock()

	dataVal, ok := a.data[fingerprint]
	if !ok {
		return fmt.Errorf("actor not found, can't set tx data")
	}

	if dataVal.txData == nil {
		dataVal.txData = make(map[string]map[string]any)
	}

	txData, ok := dataVal.txData[transactionID]
	if !ok {
		values := make(map[string]any)
		values[key] = value

		txData = make(map[string]any)
	}

	txData[key] = value
	dataVal.txData[transactionID] = txData
	dataVal.ttl = time.Now().Add(internal.ActorCacheTTL)
	a.data[fingerprint] = dataVal

	return nil
}

func (a *Store) GetActor(fingerprint string) *models.Actor {
	a.purge()
	a.lock.Lock()
	defer a.lock.Unlock()

	dataVal, ok := a.data[fingerprint]
	if !ok {
		return nil
	}

	a.data[fingerprint] = data{
		actorData: dataVal.actorData,
		ttl:       time.Now().Add(internal.ActorCacheTTL),
	}

	return dataVal.actorData
}

func (a *Store) SetActor(fingerprint string, actorData *models.Actor) {
	a.lock.Lock()
	defer a.lock.Unlock()

	ttl := time.Now().Add(internal.ActorCacheTTL)

	a.data[fingerprint] = data{
		actorData: actorData,
		ttl:       ttl,
	}
}

// purge removes actors from the store that have expired
// It also logs the size of the store to help with debugging and tracking memory usage
func (a *Store) purge() {
	a.lock.Lock()
	defer a.lock.Unlock()
	if len(a.data) > 100 { // Logging to follow up on store size
		slogctx.Info(nil, "actor store size", slog.Int("len", len(a.data)))
		if len(a.data) > 500 {
			slogctx.Warn(nil, "actor store size", slog.Int("len", len(a.data)))
			if len(a.data) > 1000 {
				slogctx.Error(nil, "actor store size", slog.Int("len", len(a.data)))
			}
		}
	}
	for fp, actor := range a.data {
		if actor.ttl.Before(time.Now()) {
			delete(a.data, fp)
		}
	}
}
