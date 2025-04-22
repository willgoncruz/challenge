package model

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StorageStorageAndPickup(t *testing.T) {
	storage := NewStorage(&sync.Mutex{}, Cold, 6)

	err := storage.Store(Order{ID: "1"})
	assert.Nil(t, err)
	assert.False(t, storage.Full())

	err = storage.Pickup(Order{ID: "1"})
	assert.Nil(t, err)
}

func Test_StorageFullReturnErrorOnNewStore(t *testing.T) {
	storage := NewStorage(&sync.Mutex{}, Hot, 1)

	_ = storage.Store(Order{ID: "1"})
	err := storage.Store(Order{ID: "1"})
	assert.Equal(t, ErrFull, err)
	assert.True(t, storage.Full())
}
