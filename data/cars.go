package data

import (
	"context"
	"fmt"
	"github.com/CassioRoos/grpc_currency/protos/currency"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		// simple case
		rates map[string]float64
		// GRPC client
		rateClient currency.Currency_SubscribeRatesClient
	}
)

func NewCarsRepository(c currency.CurrencyClient, l hclog.Logger) CarsRepositoryInterface {
	cr := &CarsRepository{c, l, make(map[string]float64), nil}
	go cr.handleUpdates()
	return cr
}

// Responsible to handle the updates and update cache
func (c *CarsRepository) handleUpdates() {
	sub, err := c.currency.SubscribeRates(context.Background())
	if err != nil {
		c.log.Error("Unable to subscribe for rates", "error", err)
		return
	}
	c.rateClient = sub
	// receive is blocking
	for {
		rrStream, err := sub.Recv()
		if grpcError := rrStream.GetError(); grpcError != nil {
			c.log.Error("Error subscribing for rates", "error", grpcError.Message)
		}
		if resp := rrStream.GetRateResponse(); resp != nil {
			if err != nil {
				c.log.Error("Error receiving message", "error", err)
				return
			}
			c.log.Info("Update received", "destination", resp.Destination.String(), "rate", resp.Rate)
			c.rates[resp.Destination.String()] = resp.Rate
		}
	}

}

// return all the cars in the DB
func (c *CarsRepository) GetCars(cur string) (Cars, error) {
	if strings.TrimSpace(cur) == "" {
		return carList, nil
	}

	rate, err := c.getRate(cur)
	if err != nil {
		c.log.Error("Unable to get rate", "Currency", cur, err)
		return nil, err
	}
	cr := Cars{}
	for _, car := range carList {
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
	rate, err := c.getRate(cur)
	if err != nil {
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

func (c *CarsRepository) getRate(destination string) (float64, error) {
	// if cached return
	if _, ok := c.rates[destination]; ok {
		return c.rates[destination], nil
	}
	rr := &currency.RateRequest{
		Base:        currency.Currencies(currency.Currencies_value["BRL"]),
		Destination: currency.Currencies(currency.Currencies_value[destination])}
	// get initial rate
	resp, err := c.currency.GetRate(context.Background(), rr)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			// I`m sure that will be the first item
			md := s.Details()[0].(*currency.RateRequest)
			if s.Code() == codes.InvalidArgument {
				return -1, fmt.Errorf(
					"Unable to get rate from currency server, base and destination currencies can not be the same base: %s, destination: %s",
					md.Base.String(),
					md.Destination.String())
			}
			return -1, fmt.Errorf(
				"Unable to get rate from currency server, base: %s, destination: %s",
				md.Base.String(),
				md.Destination.String())
		}

	}
	// subscribe for future updates
	err = c.rateClient.Send(rr)
	// set the value to cache
	c.rates[destination] = resp.Rate
	if err != nil {
		c.log.Error("Unale to subscribe to rates update", "error", err, "destination", destination)
	}
	if err != nil {
		c.log.Error("Unable to get the rate from currency", "currency", destination)
		return 1, err
	}
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
