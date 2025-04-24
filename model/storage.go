package model

import (
	"sync"
)

// interface to define all storage functionality
type Storage interface {
	Store(order Order) error
	Pickup(order Order) error
	Replace(old, new Order)
	Apply(f func(key any, value any) bool)
	Full() bool
	Empty() bool
	IsIdealTemp(temp Temperature) bool
}

type basicStorage struct {
	temp     Temperature // temperature of the Storage
	items    sync.Map    // items stored inside sync map
	count    int         // maximum capacity for the storage
	capacity int         // maximum capacity for the storage
}

func NewStorage(temp Temperature, capacity int) Storage {
	return &basicStorage{
		temp:     temp,
		items:    sync.Map{},
		capacity: capacity,
	}
}

func (s *basicStorage) Store(order Order) error {
	if s.Full() {
		return ErrFull
	}

	s.items.Store(order.ID, order)
	s.count += 1

	return nil
}

func (s *basicStorage) Pickup(order Order) error {
	_, ok := s.items.Load(order.ID)
	if !ok {
		return ErrNotFound
	}

	s.count -= 1
	return nil
}

func (s *basicStorage) Replace(oldOrder, newOrder Order) {
	s.items.Delete(oldOrder.ID)
	s.items.Store(newOrder.ID, newOrder)
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
