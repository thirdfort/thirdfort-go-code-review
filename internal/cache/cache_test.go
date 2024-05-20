package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thirdfort/thirdfort-go-code-review/internal"
	"github.com/thirdfort/thirdfort-go-code-review/internal/models"
)

func TestActor(t *testing.T) {
	fingerprint := "fingerprint"
	model := &models.Actor{ID: internal.StrPtr("actor-id"), Name: "actor-name"}

	t.Run("No data in store", func(t *testing.T) {
		store := New()
		actor := store.GetActor(fingerprint)
		assert.Nil(t, actor)
	})

	t.Run("No data in store, save data, get saved data", func(t *testing.T) {
		store := New()
		actor := store.GetActor(fingerprint)
		assert.Nil(t, actor)

		store.SetActor(fingerprint, model)
		actor = store.GetActor(fingerprint)
		assert.Equal(t, model, actor)
	})

	t.Run("Cache TTL expiration", func(t *testing.T) {
		store := New()
		actor := store.GetActor(fingerprint)
		assert.Nil(t, actor)

		store.SetActor(fingerprint, model)

		actor = store.GetActor(fingerprint)
		assert.Equal(t, model, actor)

		cache := store.data[fingerprint]
		cache.ttl = time.Now().Add(time.Duration(5 * -1 * time.Second))
		store.data[fingerprint] = cache

		actor = store.GetActor(fingerprint)
		assert.Nil(t, actor)
	})
}

func TestTransaction(t *testing.T) {
	fingerprint := "fingerprint"
	model := &models.Actor{ID: internal.StrPtr("actor-id"), Name: "actor-name"}

	t.Run("No data in store - get", func(t *testing.T) {
		store := New()
		actor := store.GetTxData(fingerprint, "transaction-id", "key")
		assert.Nil(t, actor)
	})

	t.Run("No data in store - set", func(t *testing.T) {
		store := New()
		err := store.SetTxData(fingerprint, "transaction-id", "key", "value")
		assert.NotNil(t, err)
	})

	t.Run("No tx data in store - get", func(t *testing.T) {
		store := New()

		store.SetActor(fingerprint, model)
		val := store.GetTxData(fingerprint, "transaction-id", "key")
		assert.Nil(t, val)
	})

	t.Run("No tx data in store - set and get", func(t *testing.T) {
		store := New()

		store.SetActor(fingerprint, model)

		err := store.SetTxData(fingerprint, "transaction-id", "key", "value")
		assert.Nil(t, err)

		val := store.GetTxData(fingerprint, "transaction-id", "key")
		assert.Equal(t, "value", val)
	})

	t.Run("Tx data in store - set more and get", func(t *testing.T) {
		store := New()

		store.SetActor(fingerprint, model)

		err := store.SetTxData(fingerprint, "transaction-id", "key", "value")
		assert.Nil(t, err)

		err = store.SetTxData(fingerprint, "transaction-id", "key2", "value2")
		assert.Nil(t, err)

		val := store.GetTxData(fingerprint, "transaction-id", "key")
		assert.Equal(t, "value", val)

		val2 := store.GetTxData(fingerprint, "transaction-id", "key2")
		assert.Equal(t, "value2", val2)
	})

	t.Run("Cache TTL expiration - tx", func(t *testing.T) {
		store := New()
		actor := store.GetActor(fingerprint)
		assert.Nil(t, actor)

		store.SetActor(fingerprint, model)

		actor = store.GetActor(fingerprint)
		assert.Equal(t, model, actor)

		err := store.SetTxData(fingerprint, "transaction-id", "key", "value")
		assert.Nil(t, err)

		val := store.GetTxData(fingerprint, "transaction-id", "key")
		assert.Equal(t, "value", val)

		cache := store.data[fingerprint]
		cache.ttl = time.Now().Add(time.Duration(5 * -1 * time.Second))
		store.data[fingerprint] = cache

		actor = store.GetActor(fingerprint)
		assert.Nil(t, actor)

		val = store.GetTxData(fingerprint, "transaction-id", "key")
		assert.Nil(t, val)
	})

	t.Run("Delete key from tx cache", func(t *testing.T) {
		store := New()
		actor := store.GetActor(fingerprint)
		assert.Nil(t, actor)

		store.SetActor(fingerprint, model)

		actor = store.GetActor(fingerprint)
		assert.Equal(t, model, actor)

		err := store.SetTxData(fingerprint, "transaction-id", "key", "value")
		assert.Nil(t, err)

		val := store.GetTxData(fingerprint, "transaction-id", "key")
		assert.Equal(t, "value", val)

		store.ResetTxData(fingerprint, "transaction-id", "key")

		val = store.GetTxData(fingerprint, "transaction-id", "key")
		assert.Nil(t, val)
	})
}
