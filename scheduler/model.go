package scheduler

import (
	"challenge/model"
	"sync"
)

// Interface for a async process of orders
type Scheduler interface {
	Process(orders []model.Order, wg sync.WaitGroup)
}
