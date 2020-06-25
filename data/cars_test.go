package data

import "testing"

func TestCar_Validate(t *testing.T) {
	c := &Car{
		Name:         "A",
		Price:        1.00,
		LicensePlate: "AVX-9999",
	}

	if err := c.Validate(); err != nil {
		t.Fatal(err)
	}

}
