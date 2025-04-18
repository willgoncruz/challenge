package scheduler

import (
	"challenge/kitchen"
	"challenge/ledger"
	"challenge/model"
	"log"
	"math/rand/v2"
	"sync"
	"time"
)

type PickupScheduler struct {
	min *time.Duration
	max *time.Duration
}

func NewPickupScheduler(min, max *time.Duration) *PickupScheduler {
	return &PickupScheduler{
		min: min,
		max: max,
	}
}

func (p *PickupScheduler) Process(orders []model.Order, wg *sync.WaitGroup) {
	wg.Add(1) // Add one more wait for this process

	// Wait for a random microsecond interval [min, max)
	interval := rand.Int64N(p.max.Microseconds()-p.min.Microseconds()) + p.min.Microseconds()
	time.Sleep(time.Duration(interval) * time.Microsecond)

	order := orders[0]
	if kitchen.Pickup(order) == nil {
		log.Printf("Pickup order: %+v", order)
		ledger.Audit(order, model.Pickup)
	}

	wg.Done()
}
