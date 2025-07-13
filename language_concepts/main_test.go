package main

import (
	"testing"
)

func TestProcessTruck(t *testing.T) {
	t.Run("should load and unload cargo", func(t *testing.T) {
		nt := &NormalTruck{id: "NormalTruck1", cargo: 42}
		err := processTruck(nt)
		if err != nil {
			t.Fatalf("error processing truck: %s\n", err)
		}

		et := &ElectricTruck{id: "ElectricTruck1"}
		err = processTruck(et)
		if err != nil {
			t.Fatalf("Error processing truck: %s", err)
		}

		//asserting
		if nt.cargo != 0 {
			t.Fatal("Cargo should be 0 but got \n")
		}

		if et.battery != -2 {
			t.Fatal("Battery should be -2 but got \n")
		}
	})
}
