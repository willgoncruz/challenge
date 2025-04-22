package model

import (
	"container/heap"
	"sync"
)

type ShelfStorage interface {
	DiscardCandidate() Order
}

type shelfStorage struct {
	temp     Temperature // temperature of the Storage
	items    sync.Map    // items stored inside sync map
	mtx      *sync.Mutex // Mutual exclusion
	count    int         // maximum capacity for the storage
	capacity int         // maximum capacity for the storage

	discardQueue *MinHeap // discard candidate control on min heap
}

func NewShelfStorage(mtx *sync.Mutex, capacity int) Storage {
	discardQueue := &MinHeap{}
	heap.Init(discardQueue)
	return &shelfStorage{
		mtx:          mtx,
		temp:         Room,
		items:        sync.Map{},
		capacity:     capacity,
		discardQueue: discardQueue,
	}
}

func (s *shelfStorage) Store(order Order) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

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
	s.mtx.Lock()
	defer s.mtx.Unlock()

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

func (s *shelfStorage) Remove(order Order) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.items.Delete(order.ID)
	s.count -= 1
	if pos := s.discardQueue.Find(order); pos != -1 {
		heap.Remove(s.discardQueue, pos)
	}
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
