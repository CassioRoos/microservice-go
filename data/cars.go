package data

import (
	"context"
	"fmt"
	"github.com/CassioRoos/grpc_currency/protos/currency"
	"github.com/hashicorp/go-hclog"
	"strings"
)

// Is an error raised when a car is not found
var ErrCarNotFound = fmt.Errorf("Car not found")

// Car defines the structure for an API car
// swagger:model
type (
	Car struct {
		// the id for the car
		//
		// required: false
		// min: 1
		ID int `json:"id"`

		// the color of the car
		//
		// required: false
		Color string `json:"color"`
		// the name of the car
		//
		// required: true
		// max length: 255
		Name string `json:"name" validate:"required"`
		// the description for this car
		//
		// required: false
		// max length: 255
		Description string `json:"description"`

		// the price for the car
		//
		// required: true
		// min: 0.01
		Price float64 `json:"price" validate:"required,gt=0"`

		// the license plate for this car
		//
		// required: true
		// pattern: [A-Z]{3}-[0-9]{4}
		LicensePlate string `json:"license_plate" validate:"required,lcplt"`
	}
	// This type is to help structure the code, make some changes more independent
	Cars []*Car

	CarsRepositoryInterface interface {
		GetCars(cur string) (Cars, error)
		GetCarById(id int, cur string) (*Car, error)
		UpdateCar(car Car) error
		DeleteCar(id int) error
		AddCar(car *Car)
	}

	CarsRepository struct {
		currency currency.CurrencyClient
		log      hclog.Logger
	}
)

func NewCarsRepository(c currency.CurrencyClient, l hclog.Logger) CarsRepositoryInterface {
	return &CarsRepository{c, l}
}

// return all the cars in the DB
func (c *CarsRepository) GetCars(cur string) (Cars, error) {
	if strings.TrimSpace(cur) == "" {
		return carList, nil
	}

	rate, err :=c.getRate(cur)
	if err != nil{
		c.log.Error("Unable to get rate", "Currency", cur, err)
		return nil, err
	}
	cr := Cars{}
	for _, car := range carList{
		nc := *car
		nc.Price *= rate
		cr = append(cr, &nc)

	}
	return cr, nil
}

// return a specific car by the given ID
func (c *CarsRepository) GetCarById(id int, cur string) (*Car, error) {
	i := c.FindIndexyCarId(id)
	if id == -1 {
		return nil, ErrCarNotFound
	}
	if strings.TrimSpace(cur) == "" {
		return carList[i], nil
	}
	rate, err :=c.getRate(cur)
	if err != nil{
		c.log.Error("Unable to get rate", "Currency", cur, err)
		return nil, err
	}
	nc := *carList[i]
	nc.Price *= rate
	return &nc, nil

}

// Finds the index of a car in the Database
// returns -1 when no car can be found
func (c *CarsRepository) FindIndexyCarId(id int) int {
	for i, car := range carList {
		if car.ID == id {
			return i
		}
	}
	return -1
}

// DeleteCar deletes a car from database
func (c *CarsRepository) DeleteCar(id int) error {
	i := c.FindIndexyCarId(id)
	if i == -1 {
		return ErrCarNotFound
	}
	carList = append(carList[:i], carList[i+1:]...)
	return nil
}

// AddCar adds a new car to DB
func (c *CarsRepository) AddCar(car *Car) {
	car.ID = carList[len(carList)-1].ID + 1
	carList = append(carList, car)
}

// Update a car by the given ID.
// If a car does not exist by the given id an error is returned
// CarNotFound error
func (c *CarsRepository) UpdateCar(car Car) error {
	pos := c.FindIndexyCarId(car.ID)
	if pos == -1 {
		return ErrCarNotFound
	}
	// update the car in the DB
	carList[pos] = &car
	return nil
}

func (c *CarsRepository) getRate(destination string) (float64, error){
	rr := &currency.RateRequest{
		Base: currency.Currencies(currency.Currencies_value["BRL"]),
		Destination: currency.Currencies(currency.Currencies_value[destination])}
	resp, err := c.currency.GetRate(context.Background(), rr)
	return resp.Rate, err
}


var carList = []*Car{
	&Car{ID: 1,
		Name:         "Cruze",
		Color:        "Blue",
		Description:  "A family car",
		Price:        12461.85,
		LicensePlate: "IVP-5464",
	},
	&Car{ID: 2,
		Name:         "Celta",
		Color:        "Red",
		Description:  "Economic car",
		Price:        837.37,
		LicensePlate: "ABC-4321",
	},
}

// ***** **** *** ** * Note * ** *** **** *****
// ***** **** *** ** * Note * ** *** **** *****
//
// To implement swagger look into the links
//
//
//Contents:
//Swagger Go Code Generator:
//https://github.com/go-swagger/go-swagger
//
//Swagger:
//https://swagger.io/
//
//ReDoc:
//https://github.com/Redocly/redoc
//
//Middleware for hosting redoc sites from your API:
//https://github.com/go-openapi/runtime/tree/master/middleware
//
//

// ***** **** *** ** * Note * ** *** **** *****
// ***** **** *** ** * Note * ** *** **** *****
// ***** **** *** ** * Note * ** *** **** *****
