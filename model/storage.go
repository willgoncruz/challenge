package model

import (
	"sync"
)

type Storage interface {
	Store(order Order) error
	Pickup(order Order) error
	Remove(order Order)
	Apply(f func(key any, value any) bool)
	Full() bool
	Empty() bool
	IsIdealTemp(temp Temperature) bool
}

type basicStorage struct {
	temp     Temperature // temperature of the Storage
	items    sync.Map    // items stored inside sync map
	mtx      *sync.Mutex // Mutual exclusion
	count    int         // maximum capacity for the storage
	capacity int         // maximum capacity for the storage
}

func NewStorage(mtx *sync.Mutex, temp Temperature, capacity int) Storage {
	return &basicStorage{
		mtx:      mtx,
		temp:     temp,
		items:    sync.Map{},
		capacity: capacity,
	}
}

func (s *basicStorage) Store(order Order) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.Full() {
		return ErrFull
	}

	s.items.Store(order.ID, order)
	s.count += 1

	return nil
}

func (s *basicStorage) Pickup(order Order) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	_, ok := s.items.Load(order.ID)
	if !ok {
		return ErrNotFound
	}

	s.count -= 1
	return nil
}

func (s *basicStorage) Remove(order Order) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.items.Delete(order.ID)
	s.count -= 1
}

// Apply a function to all storage items
func (s *basicStorage) Apply(f func(key any, value any) bool) {
	s.items.Range(f)
}

func (s *basicStorage) Full() bool {
	return s.count >= s.capacity
}

func (s *basicStorage) Empty() bool {
	return s.count == 0
}

func (s *basicStorage) IsIdealTemp(temp Temperature) bool {
	return s.temp == temp
}
