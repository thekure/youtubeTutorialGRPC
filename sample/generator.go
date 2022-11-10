package sample

import (
	ytTut "github.com/thekure/youtubeTutorialGRPC/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// NewKeyboard returns a new sample keyboard
func NewKeyboard() *ytTut.Keyboard {
	keyboard := &ytTut.Keyboard{
		Layout:  randomKeyboardLayout(),
		Backlit: randomBool(),
	}
	return keyboard
}

// Returns a new sample CPU
func NewCPU() *ytTut.CPU {
	brand := randomCPUBrand()
	name := randomCPUName(brand)

	numberCores := randomInt(2, 8)
	numberThreads := randomInt(numberCores, 12)

	minGhz := randomFloat64(2.0, 3.5)
	maxGhz := randomFloat64(minGhz, 5.0)

	cpu := &ytTut.CPU{
		Brand:         brand,
		Name:          name,
		NumberCores:   uint32(numberCores),
		NumberThreads: uint32(numberThreads),
		MinGhz:        minGhz,
		MaxGhz:        maxGhz,
	}
	return cpu
}

// Returns a new sample GPU
func newGPU() *ytTut.GPU {
	brand := randomGPUBrand()
	name := randomGPUName(brand)

	minGhz := randomFloat64(1.0, 1.5)
	maxGhz := randomFloat64(minGhz, 2.0)

	memory := &ytTut.Memory{
		Value: uint64(randomInt(2, 6)),
		Unit:  ytTut.Memory_GIGABYTE,
	}

	gpu := &ytTut.GPU{
		Brand:  brand,
		Name:   name,
		MinGhz: minGhz,
		MaxGhz: maxGhz,
		Memory: memory,
	}
	return gpu
}

// Returns a new sample RAM
func NewRam() *ytTut.Memory {
	ram := &ytTut.Memory{
		Value: uint64(randomInt(4, 64)),
		Unit:  ytTut.Memory_GIGABYTE,
	}
	return ram
}

// Returns a new sample SSD storage
func NewSSD() *ytTut.Storage {
	ssd := &ytTut.Storage{
		Driver: ytTut.Storage_SDD,
		Memory: &ytTut.Memory{
			Value: uint64(randomInt(128, 1024)),
			Unit:  ytTut.Memory_GIGABYTE,
		},
	}

	return ssd
}

// Returns a new sample HDD storage
func NewHDD() *ytTut.Storage {
	hdd := &ytTut.Storage{
		Driver: ytTut.Storage_HDD,
		Memory: &ytTut.Memory{
			Value: uint64(randomInt(1, 6)),
			Unit:  ytTut.Memory_TERABYTE,
		},
	}
	return hdd
}

// Returns new sample screen
func NewScreen() *ytTut.Screen {

	screen := &ytTut.Screen{
		SizeInch:   float32(randomFloat64(13, 17)),
		Resolution: randomScreenResolution(),
		Panel:      randomScreenPanel(),
		Multitouch: randomBool(),
	}
	return screen
}

// Return new sample Laptop
func NewLaptop() *ytTut.Laptop {
	brand := randomLaptopBrand()
	name := randomLaptopName(brand)

	laptop := &ytTut.Laptop{
		Id:          randomID(),
		Brand:       brand,
		Name:        name,
		Cpu:         NewCPU(),
		Ram:         NewRam(),
		Gpus:        []*ytTut.GPU{newGPU()},
		Stoages:     []*ytTut.Storage{NewSSD(), NewHDD()},
		Screen:      NewScreen(),
		Keyboard:    NewKeyboard(),
		Weight:      &ytTut.Laptop_WeightKg{WeightKg: randomFloat64(1.0, 3.0)},
		PriceUsd:    randomFloat64(1500, 3000),
		ReleaseYear: uint32(randomInt(2015, 2019)),
		UpdatedAt:   timestamppb.Now(),
	}

	return laptop
}
