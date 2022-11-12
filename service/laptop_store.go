package service

import (
	"errors"
	"fmt"
	"sync"

	"github.com/jinzhu/copier"
	ytTut "github.com/thekure/youtubeTutorialGRPC/grpc"
)

var ErrAlreadyExists = errors.New("record already exists")

// Interface to store Laptops in memory
type LaptopStore interface {
	// Saves laptop to store
	Save(laptop *ytTut.Laptop) error
	// Finds laptop by ID:
	Find(id string) (*ytTut.Laptop, error)
}

type InMemoryLaptopStore struct {
	mutex sync.RWMutex
	data  map[string]*ytTut.Laptop
}

// Returns new InMemoryLaptopStore
func NewInMemoryLaptopStore() *InMemoryLaptopStore {
	return &InMemoryLaptopStore{
		data: make(map[string]*ytTut.Laptop),
	}
}

func (store *InMemoryLaptopStore) Save(laptop *ytTut.Laptop) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.data[laptop.Id] != nil {
		return ErrAlreadyExists
	}

	// deep copy of laptop object (requires: go get github.com/jinzhu/copier)
	// copiers the laptop object in other, and saves other in store, by its Id.
	other := &ytTut.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return fmt.Errorf("cannot copy laptop data: %w", err)
	}

	store.data[other.Id] = other
	return nil
}

func (store *InMemoryLaptopStore) Find(id string) (*ytTut.Laptop, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	laptop := store.data[id]
	if laptop == nil {
		return nil, nil
	}

	other := &ytTut.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, fmt.Errorf("cannot copy laptop data: %w", err)
	}

	return other, nil
}
