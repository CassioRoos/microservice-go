package data

import (
	"encoding/json"
	"io"
	"time"
)

type (
	Car struct {
		ID           int     `json:"id"`
		Color        string  `json:"color"`
		Name         string  `json:"name"`
		Description  string  `json:"description"`
		Price        float32 `json:"price"`
		LicensePlate string  `json:"license_plate"`
		CreatedOn    string  `json:"-"`
		UpdatedOn    string  `json:"-"`
		DeletedOn    string  `json:"-"`
	}
	// This type is to help structure the code, make some changes more independent
	Cars []*Car
)

func (c *Car) FromJSON(r io.Reader) error {
	// the opposite of encode
	d := json.NewDecoder(r)
	return d.Decode(c)
}

// This works like a object method
func (c *Cars) ToJSON(w io.Writer) error {
	// User the encoder because is faster to marshal json | Marshal uses in memory stream to perform
	e := json.NewEncoder(w)
	// now we have to encode ourselfs, because C is pointing to the slice of cars
	return e.Encode(c)
}

func GetCars() Cars {
	return carList
}

var carList = []*Car{
	&Car{ID: 1,
		Name:         "Cruze",
		Color:        "Blue",
		Description:  "A family car",
		Price:        12461.85,
		LicensePlate: "IVP5464",
		CreatedOn:    time.Now().UTC().String(),
		UpdatedOn:    time.Now().UTC().String(),
	},
	&Car{ID: 2,
		Name:         "Celta",
		Color:        "Red",
		Description:  "Economic car",
		Price:        837.37,
		LicensePlate: "X897H97",
		CreatedOn:    time.Now().UTC().String(),
		UpdatedOn:    time.Now().UTC().String(),
	},
}

func AddCar(c *Car) {
	c.ID = getNext()
	carList = append(carList, c)
}

func getNext() int {
	c := carList[len(carList)-1]
	return c.ID + 1
}
