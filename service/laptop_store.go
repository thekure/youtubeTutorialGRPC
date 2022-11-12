package service

import (
	"context"
	"errors"
	"fmt"
	"log"
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
	// Searches for laptops with filter. Returns one by one via the found function (callback)
	Search(ctx context.Context, filter *ytTut.Filter, found func(laptop *ytTut.Laptop) error) error
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
	other, err := deepCopy(laptop)

	if err != nil {
		return err
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

	return deepCopy(laptop)
}

func (store *InMemoryLaptopStore) Search(
	ctx context.Context,
	filter *ytTut.Filter,
	found func(laptop *ytTut.Laptop) error,
) error {
	// since we are reading data, we have to use a lock:
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	for _, laptop := range store.data {

		if ctx.Err() == context.Canceled || ctx.Err() == context.DeadlineExceeded {
			log.Print("context is cancelled")
			return errors.New("context is cancelled")
		}

		if isQualified(filter, laptop) {
			other, err := deepCopy(laptop)
			if err != nil {
				return err
			}

			found(other)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func isQualified(filter *ytTut.Filter, laptop *ytTut.Laptop) bool {
	if laptop.GetPriceUsd() > filter.GetMaxPriceUsd() {
		return false
	}

	if laptop.GetCpu().GetNumberCores() < filter.GetMinCpuCores() {
		return false
	}

	if laptop.GetCpu().GetMinGhz() < filter.GetMinCpuGhz() {
		return false
	}

	if toBit(laptop.GetRam()) < toBit(filter.GetMinRam()) {
		return false
	}

	return true
}

func toBit(memory *ytTut.Memory) uint64 {
	value := memory.GetValue()

	switch memory.GetUnit() {
	case ytTut.Memory_BIT:
		return value
	case ytTut.Memory_BYTE:
		return value * 8 // or value << 3 because 8 = 2^3
	case ytTut.Memory_KILOBYTE:
		return value << 13
	case ytTut.Memory_MEGABYTE:
		return value << 23
	case ytTut.Memory_GIGABYTE:
		return value << 33
	case ytTut.Memory_TERABYTE:
		return value << 43
	default:
		return 0
	}
}

func deepCopy(laptop *ytTut.Laptop) (*ytTut.Laptop, error) {
	other := &ytTut.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, fmt.Errorf("cannot copy laptop data: %w", err)
	}

	return other, nil
}
