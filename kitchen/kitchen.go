package kitchen

import (
	"challenge/ledger"
	"challenge/model"
	"sync"
)

type Kitchen struct {
	cooler model.Storage
	heater model.Storage
	shelf  model.Storage
}

var kitchen *Kitchen

func init() {
	reset()
}

func reset() {
	mtx := &sync.Mutex{}
	kitchen = &Kitchen{
		cooler: model.NewStorage(mtx, model.Cold, 6),
		heater: model.NewStorage(mtx, model.Hot, 6),
		shelf:  model.NewShelfStorage(mtx, 12),
	}
}

func Place(order model.Order) {
	// Decision put the order on the ideal storage
	idealStorage := getIdealStorageByTemp(order.Temp)
	if idealStorage.Store(order) == nil {
		return // Was added to ideal storage
	}

	// Otherwise, must put order on the shelf first
	if kitchen.shelf.Store(order) == nil {
		return // Was added on shelf
	}

	// Shelf is full, try to move something from shelf to other storages
	moved := false
	kitchen.shelf.Apply(func(key any, value any) bool {
		shelfOrder := value.(model.Order)
		idealStorage := getIdealStorageByTemp(shelfOrder.Temp)
		if idealStorage.Store(shelfOrder) == nil { // Could move something from the shelf to ideal storage
			kitchen.shelf.Replace(shelfOrder, order)
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
	discardCandidate := kitchen.shelf.(model.ShelfStorage).DiscardCandidate()

	kitchen.shelf.Replace(discardCandidate, order) // replace discard with placed order
	ledger.Audit(discardCandidate, model.Discard)  // audit the discard move
}

func Pickup(order model.Order) error {
	// First pickup on shelf
	err := kitchen.shelf.Pickup(order)
	if err == nil {
		return nil
	}

	// Then on ideal storage
	return getIdealStorageByTemp(order.Temp).Pickup(order)
}

func getIdealStorageByTemp(temp model.Temperature) model.Storage {
	if kitchen.cooler.IsIdealTemp(temp) {
		return kitchen.cooler
	}

	if kitchen.heater.IsIdealTemp(temp) {
		return kitchen.heater
	}

	return kitchen.shelf
}
