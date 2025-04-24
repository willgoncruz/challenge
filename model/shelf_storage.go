package model

import (
	"container/heap"
	"sync"
)

type ShelfStorage interface {
	DiscardCandidate() Order
}

// ShelfStorage is a more complex storage as it keeps track of discard heap
type shelfStorage struct {
	temp     Temperature // temperature of the Storage
	items    sync.Map    // items stored inside sync map
	count    int         // maximum capacity for the storage
	capacity int         // maximum capacity for the storage

	discardQueue *MinHeap // discard candidate control on min heap
}

func NewShelfStorage(capacity int) Storage {
	discardQueue := &MinHeap{}
	heap.Init(discardQueue)
	return &shelfStorage{
		temp:         Room,
		items:        sync.Map{},
		capacity:     capacity,
		discardQueue: discardQueue,
	}
}

func (s *shelfStorage) Store(order Order) error {
	if s.Full() {
		return ErrFull
	}

	s.items.Store(order.ID, order)
	s.count += 1

	// Control the discard queue data
	order.FillTTL(s)
	heap.Push(s.discardQueue, order)

	return nil
}

func (s *shelfStorage) Pickup(order Order) error {
	_, ok := s.items.Load(order.ID)
	if !ok {
		return ErrNotFound
	}

	s.count -= 1
	if pos := s.discardQueue.Find(order); pos != -1 {
		heap.Remove(s.discardQueue, pos)
	}

	return nil
}

func (s *shelfStorage) Replace(oldOrder, newOrder Order) {
	s.items.Delete(oldOrder.ID)
	if pos := s.discardQueue.Find(oldOrder); pos != -1 {
		heap.Remove(s.discardQueue, pos)
	}

	s.items.Store(newOrder.ID, newOrder)
	heap.Push(s.discardQueue, newOrder)
}

// Apply a function to all storage items
func (s *shelfStorage) Apply(f func(key any, value any) bool) {
	s.items.Range(f)
}

func (s *shelfStorage) Full() bool {
	return s.count >= s.capacity
}

func (s *shelfStorage) Empty() bool {
	return s.count == 0
}

func (s *shelfStorage) IsIdealTemp(temp Temperature) bool {
	return s.temp == temp
}

func (s *shelfStorage) DiscardCandidate() Order {
	return heap.Pop(s.discardQueue).(Order)
}
