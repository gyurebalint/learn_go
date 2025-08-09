package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

var (
	ErrNotImplemented = errors.New("not implemented")
	ErrTruckNotFound  = errors.New("truck not found")
)

type Truck interface {
	LoadCargo() error
	UnloadCargo() error
}

type ElectricTruck struct {
	id      string
	cargo   int
	battery float64
}

func (e *ElectricTruck) LoadCargo() error {
	e.cargo += 1
	e.battery -= 1
	return nil
}

func (e *ElectricTruck) UnloadCargo() error {
	e.cargo += 0
	e.battery -= 1
	return nil
}

type NormalTruck struct {
	id    string
	cargo int
}

func (t *NormalTruck) LoadCargo() error {
	t.cargo += 1
	return nil
}

func (t *NormalTruck) UnloadCargo() error {
	t.cargo = 0
	return nil
}

type contextKey string

var UserIdKey contextKey = "userID"

func processTruck(ctx context.Context, truck Truck) error {
	fmt.Printf("processing truck %+v\n", truck)

	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	//simulate long-running task
	delay := time.Second * 3
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(delay):
		break
	}

	userID := ctx.Value(UserIdKey)
	log.Println(userID)

	if err := truck.LoadCargo(); err != nil {
		return fmt.Errorf("error loading cargo: %w", err)
	}
	if err := truck.UnloadCargo(); err != nil {
		return fmt.Errorf("error unloading cargo: %w", err)
	}

	fmt.Printf("finished truck %+v\n", truck)
	return nil
}

func processFleet(ctx context.Context, trucks []Truck) error {
	var wg sync.WaitGroup
	errorChannel := make(chan error, len(trucks))

	for _, t := range trucks {
		wg.Add(1)

		go func(t Truck) {
			if err := processTruck(ctx, t); err != nil {

				errorChannel <- err
			}
			wg.Done()
		}(t)
	}
	wg.Wait()
	close(errorChannel)

	var errs []error
	for err := range errorChannel {
		log.Printf("error processing truck %+v\n", err)
		errs = append(errs, err)
	}

	if len(errs) > 0 {

		return fmt.Errorf("fleet processing had %d errors", len(errs))
	}

	return nil
}

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, UserIdKey, 42)

	fleet := []Truck{
		&ElectricTruck{id: "ET1", cargo: 0, battery: 100},
		&NormalTruck{id: "NT1", cargo: 0},
		&ElectricTruck{id: "ET2", cargo: 0, battery: 100},
		&NormalTruck{id: "NT2", cargo: 0},
		&ElectricTruck{id: "ET1", cargo: 0, battery: 100},
		&NormalTruck{id: "NT1", cargo: 0},
		&ElectricTruck{id: "ET2", cargo: 0, battery: 100},
		&NormalTruck{id: "NT2", cargo: 0},
		&ElectricTruck{id: "ET1", cargo: 0, battery: 100},
		&NormalTruck{id: "NT1", cargo: 0},
		&ElectricTruck{id: "ET2", cargo: 0, battery: 100},
		&NormalTruck{id: "NT2", cargo: 0},
		&ElectricTruck{id: "ET1", cargo: 0, battery: 100},
		&NormalTruck{id: "NT1", cargo: 0},
		&ElectricTruck{id: "ET2", cargo: 0, battery: 100},
		&NormalTruck{id: "NT2", cargo: 0},
		&ElectricTruck{id: "ET1", cargo: 0, battery: 100},
		&NormalTruck{id: "NT1", cargo: 0},
		&ElectricTruck{id: "ET2", cargo: 0, battery: 100},
		&NormalTruck{id: "NT2", cargo: 0},
	}

	if err := processFleet(ctx, fleet); err != nil {
		fmt.Printf("error processing fleet: %v\n", err)
		return
	}

	fmt.Println("All trucks processed successfully")
}

/*func main() {
    nt := &NormalTruck{id: "NormalTruck1"}
    err := processTruck(nt)
    if err != nil {
       log.Fatalf("error processing truck: %s\n", err)
    }

    et := &ElectricTruck{id: "ElectricTruck1"}
    err = processTruck(et)
    if err != nil {
       log.Fatalf("Error processing truck: %s", err)
    }

    log.Println(nt.cargo)
    log.Println(et.battery)
}*/
