package scheduler

import (
	"challenge/kitchen"
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
		log.Printf("Placed order: %+v", order)

		// Async start the pickup process after placing order
		wg.Add(1)
		go p.pickupScheduler.Process(order, wg)

		time.Sleep(*p.rate)
	}

	wg.Done()
}
