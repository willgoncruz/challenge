package kitchen

import (
	"challenge/model"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Complex test with full storage, should get the closest to spoiling
func Test_DiscardCandidateShouldBeClosestToSpoiled(t *testing.T) {
	for i := range 6 { // fill cooler and heater
		Place(CommonOrder(strconv.Itoa(i)+"-cold", model.Cold))
		Place(CommonOrder(strconv.Itoa(2*i)+"-hot", model.Hot))
	}

	assert.True(t, kitchen.cooler.Full())
	assert.True(t, kitchen.heater.Full())
	assert.True(t, kitchen.shelf.Empty())

	spoiled := SpoiledOrder("spoiled", model.Cold)
	Place(spoiled)
	assert.False(t, kitchen.shelf.Empty())

	for i := range 11 { // fill rest of shelf
		Place(CommonOrder(strconv.Itoa(3*i)+"-shelf", model.Cold))
	}

	assert.True(t, kitchen.shelf.Full())

	// Place last order to verify discard of spoiled
	Place(CommonOrder(strconv.Itoa(3*12), model.Cold))

	// Spoiler should not be in the shelf anymore
	err := kitchen.shelf.Pickup(spoiled)
	assert.Equal(t, model.ErrNotFound, err)
}

func CommonOrder(id string, temp model.Temperature) model.Order {
	return model.Order{
		ID:        id,
		Name:      "test order",
		Temp:      temp,
		Freshness: 120,
	}
}

func SpoiledOrder(id string, temp model.Temperature) model.Order {
	return model.Order{
		ID:        id,
		Name:      "spoiled order",
		Temp:      temp,
		Freshness: 1, // The least fresh one
	}
}
