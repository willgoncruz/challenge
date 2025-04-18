package model

import (
	"sync"
)

type Storage struct {
	temp     Temperature // temperature of the Storage
	items    sync.Map    // items stored inside sync map
	mtx      *sync.Mutex // Mutual exclusion
	count    int         // maximum capacity for the storage
	capacity int         // maximum capacity for the storage
}

func NewStorage(temp Temperature, capacity int) *Storage {
	return &Storage{
		temp:     temp,
		items:    sync.Map{},
		capacity: capacity,
		mtx:      &sync.Mutex{},
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
	return nil
}

func (s *Storage) Remove(order Order) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.items.Delete(order.ID)
	s.count -= 1
}

// Apply a function to all storage items
func (s *Storage) Apply(f func(key any, value any) bool) {
	s.items.Range(f)
}

func (s *Storage) Full() bool {
	return s.count >= s.capacity
}

func (s *Storage) IsShelf() bool {
	return s.temp == Room
}
