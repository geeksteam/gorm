package grom

import (
	"errors"
	"sync"
)

const (
	MaxRoutines = 2
)

type Routine interface {
	Start()            // Starting routin body and writes info about start/stop
	OnError(err error) // Calling when panic occurred
	String() string    // Must return full info about itself
	Finished() bool
}

var (
	mutex       sync.Mutex
	routineStor = []Routine{}
)

func Run(routine Routine) {
	defer cleanHistory()

	mutex.Lock()

	routineStor = append(routineStor, routine)
	mutex.Unlock()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				switch err := r.(type) {
				case error:
					routine.OnError(err)

				case string:
					routine.OnError(errors.New(err))
				}

			}
		}()

		routine.Start()
	}()
}

// Routines returns all goroutines, which are created with Run function.
func Routines() []Routine {
	result := []Routine{}
	mutex.Lock()
	for _, v := range routineStor {
		result = append(result, v)
	}
	mutex.Unlock()

	return result
}

func cleanHistory() {
	mutex.Lock()
	defer mutex.Unlock()

	// Counting how much finished routines are stored
	var inStor = 0
	for _, routine := range routineStor {
		if routine.Finished() {
			inStor++
		}
	}

	// Check if cleaning needed
	if inStor <= MaxRoutines {
		return
	}

	// Removing from the beginning of the collection
	var offset = 0
	for i := range routineStor {
		actualIndex := i - offset

		if routineStor[actualIndex].Finished() {
			routineStor = append(routineStor[:actualIndex], routineStor[actualIndex+1:]...)
			offset++
			inStor--
		}

		if inStor <= MaxRoutines {
			break
		}
	}
}
