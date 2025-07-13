package main

import (
	"errors"
	"fmt"
	"log"
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
	e.battery += -1
	return nil
}

func (e *ElectricTruck) UnloadCargo() error {
	e.cargo += 0
	e.battery += -1
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

func processTruck(truck Truck) error {
	fmt.Printf("processing truck %+v\n", truck)

	if err := truck.LoadCargo(); err != nil {
		return fmt.Errorf("error loading cargo: %w", err)
	}

	if err := truck.UnloadCargo(); err != nil {
		return fmt.Errorf("error unloading cargo: %w", err)
	}

	return nil
}

func main() {
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
}
