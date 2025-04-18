package kitchen

import (
	"challenge/model"
	"container/heap"
	"sync"
	"time"
)

type OrderTracker struct {
	itemsByStorage sync.Map
	discardQueue   *model.MinHeap
}

func NewOrderTracker() *OrderTracker {
	discardQueue := &model.MinHeap{}
	heap.Init(discardQueue)
	return &OrderTracker{
		discardQueue:   discardQueue,
		itemsByStorage: sync.Map{},
	}
}

func (o *OrderTracker) Track(order model.Order, storage *model.Storage) {
	o.Untrack(order)

	// Track the freshness of the order
	freshnessTTL := time.Duration(order.FreshnessInSecondsByStorage(storage)) * time.Second
	order.TTL = time.Now().UnixMicro() + freshnessTTL.Microseconds()

	o.itemsByStorage.Swap(order.ID, storage)
	if storage.IsShelf() { // Track all shelf items on heap
		heap.Push(o.discardQueue, order)
	}
}

func (o *OrderTracker) Get(order model.Order) (*model.Storage, error) {
	storage, ok := o.itemsByStorage.Load(order.ID)
	if !ok {
		return nil, model.ErrNotFound
	}

	if pos := o.discardQueue.Find(order); pos != -1 {
		heap.Remove(o.discardQueue, pos)
		// heap.Fix(o.discardQueue, pos) // TODO dont know if needed
	}

	return storage.(*model.Storage), nil
}

func (o *OrderTracker) Untrack(order model.Order) {
	o.itemsByStorage.Delete(order.ID)
	if pos := o.discardQueue.Find(order); pos != -1 {
		heap.Remove(o.discardQueue, pos)
		// heap.Fix(o.discardQueue, pos) // TODO dont know if needed
	}
}

func (o *OrderTracker) DiscardCandidate() model.Order {
	return heap.Pop(o.discardQueue).(model.Order)
}
