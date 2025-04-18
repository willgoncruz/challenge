package scheduler

import (
	"challenge/kitchen"
	"challenge/ledger"
	"challenge/model"
	"log"
	"sync"
	"time"
)

type PlaceScheduler struct {
	rate            *time.Duration
	pickupScheduler *PickupScheduler
}

func NewPlaceScheduler(rate *time.Duration, pickupScheduler *PickupScheduler) *PlaceScheduler {
	return &PlaceScheduler{
		rate:            rate,
		pickupScheduler: pickupScheduler,
	}
}

func (p *PlaceScheduler) Process(orders []model.Order, wg *sync.WaitGroup) {
	for _, order := range orders {
		kitchen.Place(order)
		ledger.Audit(order, model.Place)
		log.Printf("Placed order: %+v", order)

		// Async start the pickup process after placing order
		go p.pickupScheduler.Process([]model.Order{order}, wg)

		time.Sleep(*p.rate)
	}

	wg.Done()
}
