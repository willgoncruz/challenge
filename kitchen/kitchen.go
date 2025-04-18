package kitchen

import (
	"challenge/ledger"
	"challenge/model"
	"fmt"
)

type Kitchen struct {
	cooler *model.Storage
	heater *model.Storage
	shelf  *model.Storage

	orderTracker *OrderTracker // Data structure for order tracking in performance oriented structures
}

var kitchen *Kitchen

func init() {
	kitchen = &Kitchen{
		cooler:       model.NewStorage(model.Cold, 6),
		heater:       model.NewStorage(model.Hot, 6),
		shelf:        model.NewStorage(model.Room, 12),
		orderTracker: NewOrderTracker(),
	}
}

func Place(order model.Order) {
	// Decision put the order on the ideal storage
	idealStorage := getIdealStorageByTemp(order.Temp)
	if idealStorage.Store(order) == nil {
		kitchen.orderTracker.Track(order, idealStorage)
		return // Was added to ideal storage
	}

	// Otherwise, must put order on the shelf first
	if kitchen.shelf.Store(order) == nil {
		kitchen.orderTracker.Track(order, kitchen.shelf)
		return // Was added on shelf
	}

	// Shelf is full, try to move something from shelf to other storages
	moved := false
	kitchen.shelf.Apply(func(key any, value any) bool {
		shelfOrder := value.(model.Order)
		idealStorage := getIdealStorageByTemp(shelfOrder.Temp)
		if idealStorage.Store(order) == nil { // Could move something from the shelf to ideal storage
			kitchen.shelf.Remove(order)
			kitchen.orderTracker.Track(order, idealStorage)
			ledger.Audit(shelfOrder, model.Move) // audit the move for solution
			moved = true
			return false
		}

		return true
	})

	if moved {
		return
	}

	// Last resort, must discard something
	for true {
		discardCandidate := kitchen.orderTracker.DiscardCandidate()
		storage, err := kitchen.orderTracker.Get(discardCandidate)
		if err != nil {
			fmt.Printf("Discard candidate not found on tracker: %+v", discardCandidate)
			continue // In the the discard candidate is not found, must continue the search
		}

		storage.Remove(discardCandidate)
		ledger.Audit(discardCandidate, model.Discard) // audit the discard move

		kitchen.shelf.Store(order) // add placed order to shelf after all
		kitchen.orderTracker.Track(order, kitchen.shelf)
		break // break after sucessfull discard
	}
}

func Pickup(order model.Order) error {
	storage, err := kitchen.orderTracker.Get(order)
	if err != nil {
		fmt.Printf("Order not found on any storage: %+v", order)
		return model.ErrNotFound
	}

	return storage.Pickup(order)
}

func getIdealStorageByTemp(temp model.Temperature) *model.Storage {
	if temp == model.Cold {
		return kitchen.cooler
	}

	if temp == model.Hot {
		return kitchen.heater
	}

	return kitchen.shelf
}
