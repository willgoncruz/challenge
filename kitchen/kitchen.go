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

	mtx *sync.Mutex // Mutual exclusion for concurrency handling
}

var kitchen *Kitchen

func init() {
	reset()
}

func reset() {
	mtx := &sync.Mutex{}
	kitchen = &Kitchen{
		mtx:    mtx,
		cooler: model.NewStorage(model.Cold, 6),
		heater: model.NewStorage(model.Hot, 6),
		shelf:  model.NewShelfStorage(12),
	}
}

func Place(order model.Order) {
	kitchen.mtx.Lock()
	defer kitchen.mtx.Unlock()

	// Decision put the order on the ideal storage
	idealStorage := getIdealStorageByTemp(order.Temp)
	if idealStorage.Store(order) == nil {
		ledger.Audit(order, model.Place)
		return // Was added to ideal storage
	}

	// Otherwise, must put order on the shelf first
	if kitchen.shelf.Store(order) == nil {
		ledger.Audit(order, model.Place)
		return // Was added on shelf
	}

	// Shelf is full, try to move something from shelf to other storages
	moved := false
	kitchen.shelf.Apply(func(key any, value any) bool {
		shelfOrder := value.(model.Order)
		idealStorage := getIdealStorageByTemp(shelfOrder.Temp)
		if idealStorage.Store(shelfOrder) == nil { // Could move something from the shelf to ideal storage
			kitchen.shelf.Replace(shelfOrder, order)
			ledger.Audit(shelfOrder, model.Move) // audit the move of shelf order and place of new order
			ledger.Audit(order, model.Place)
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
	ledger.Audit(discardCandidate, model.Discard)  // audit the discard move and store of new order
	ledger.Audit(order, model.Place)
}

func DiscardLeftovers() {
	// Define a discard function applying in all storages
	discardFunc := func(key any, value any) bool {
		order := value.(model.Order)
		ledger.Audit(order, model.Discard)
		return true
	}

	kitchen.cooler.Apply(discardFunc)
	kitchen.heater.Apply(discardFunc)
	kitchen.shelf.Apply(discardFunc)
}

func Pickup(order model.Order) error {
	kitchen.mtx.Lock()
	defer kitchen.mtx.Unlock()

	// First pickup on shelf
	err := kitchen.shelf.Pickup(order)
	if err == nil {
		ledger.Audit(order, model.Pickup)
		return nil
	}

	// Then on ideal storage
	err = getIdealStorageByTemp(order.Temp).Pickup(order)
	if err == nil {
		ledger.Audit(order, model.Pickup)
	}

	return err
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
