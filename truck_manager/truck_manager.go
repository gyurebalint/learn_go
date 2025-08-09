package main

import (
	"errors"
)

var ErrTruckNotFound = errors.New("truck not found")

type FleetManager interface {
	AddTruck(id string, cargo int) error
	GetTruck(id string) (*Truck, error)
	RemoveTruck(id string) error
	UpdateTruckCargo(id string, cargo int) error
}

func (tm *truckManager) AddTruck(id string, cargo int) error {
	newTruck := &Truck{ID: id, Cargo: cargo}
	tm.trucks[id] = newTruck
	return nil
}

func (tm *truckManager) GetTruck(id string) (*Truck, error) {
	var truck, exists = tm.trucks[id]
	if !exists {
		return nil, ErrTruckNotFound
	}
	return truck, nil
}

func (tm *truckManager) RemoveTruck(id string) error {
	var _, exists = tm.trucks[id]
	if !exists {
		return ErrTruckNotFound
	}
	delete(tm.trucks, id)
	return nil
}

func (tm *truckManager) UpdateTruckCargo(id string, cargo int) error {
	var truck, exists = tm.trucks[id]
	if !exists {
		return ErrTruckNotFound
	}
	truck.Cargo = cargo
	tm.trucks[id] = truck
	return nil
}

type Truck struct {
	ID    string
	Cargo int
}

type truckManager struct {
	trucks map[string]*Truck
}

func NewTruckManager() truckManager {
	return truckManager{
		trucks: make(map[string]*Truck),
	}
}
