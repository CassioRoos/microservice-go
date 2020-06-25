package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator"
)

type (
	Car struct {
		ID           int     `json:"id"`
		Color        string  `json:"color"`
		Name         string  `json:"name" validate:"required"`
		Description  string  `json:"description"`
		Price        float32 `json:"price" validate:"gt=0"`
		LicensePlate string  `json:"license_plate" validate:"required,lcplt"`
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

// ToJSON serializes the contents of the collection to JSON
// NewEncoder provides better performance than json.Unmarshal as it does not
// have to buffer the output into an in memory slice of bytes
// this reduces allocations and the overheads of the service
//
// https://golang.org/pkg/encoding/json/#NewEncoder
func (c *Cars) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	// now we have to encode ourselfs, because C is pointing to the slice of cars
	return e.Encode(c)
}

func (c *Car) Validate() error {
	v := validator.New()
	// adding a custom validation
	v.RegisterValidation("lcplt", validateLicensePlate)

	return v.Struct(c)
}

func validateLicensePlate(fl validator.FieldLevel) bool {
	// here em Brasil the license plate should be something like ABC-1234
	// regex will ensure the this is respected
	rgx := regexp.MustCompile(`[A-Z]{3}-[0-9]{4}`)
	matches := rgx.FindAllString(fl.Field().String(), -1)
	// if there is more than one or no one should raise error
	if len(matches) != 1 {
		return false
	}
	return true
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

var ErrCarNotFound = fmt.Errorf("Car not found")

func findCar(id int) (*Car, int, error) {
	for i, c := range carList {
		if c.ID == id {
			return c, i, nil
		}
	}
	return nil, -1, ErrCarNotFound
}

func UpdateCar(id int, c *Car) error {
	_, pos, err := findCar(id)
	if err != nil {
		return err
	}

	c.ID = id
	carList[pos] = c
	return nil
}
