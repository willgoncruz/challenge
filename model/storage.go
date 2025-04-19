package model

import (
	"container/heap"
	"sync"
)

type Storage struct {
	temp     Temperature // temperature of the Storage
	items    sync.Map    // items stored inside sync map
	mtx      *sync.Mutex // Mutual exclusion
	count    int         // maximum capacity for the storage
	capacity int         // maximum capacity for the storage

	discardQueue *MinHeap // discard candidate control on min heap
}

func NewStorage(temp Temperature, capacity int) *Storage {
	discardQueue := &MinHeap{}
	heap.Init(discardQueue)
	return &Storage{
		temp:         temp,
		items:        sync.Map{},
		capacity:     capacity,
		mtx:          &sync.Mutex{},
		discardQueue: discardQueue,
	}
}

func (s *Storage) Store(order Order) error {
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

func (s *Storage) Pickup(order Order) error {
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

func (s *Storage) Remove(order Order) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.items.Delete(order.ID)
	s.count -= 1
	if pos := s.discardQueue.Find(order); pos != -1 {
		heap.Remove(s.discardQueue, pos)
	}
}

// Apply a function to all storage items
func (s *Storage) Apply(f func(key any, value any) bool) {
	s.items.Range(f)
}

func (s *Storage) Full() bool {
	return s.count >= s.capacity
}

func (s *Storage) Empty() bool {
	return s.count == 0
}

func (s *Storage) IsShelf() bool {
	return s.IsIdealTemp(Room)
}

func (s *Storage) IsIdealTemp(temp Temperature) bool {
	return s.temp == temp
}

func (s *Storage) DiscardCandidate() Order {
	return heap.Pop(s.discardQueue).(Order)
}
