package handlers

import (
	"fmt"
	"github.com/CassioRoos/MicroseService/data"
	"github.com/hashicorp/go-hclog"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Cars struct {
	l  hclog.Logger
	v  *data.Validation
	cr data.CarsRepositoryInterface
}

type KeyCar struct{}

// A list of cars returns in the response
// swagger:response carsResponse
type carsResponse struct {
	// All cars in the system
	// in: body
	Body []data.Car
}

func NewCars(l hclog.Logger, v *data.Validation, cr data.CarsRepositoryInterface) *Cars {
	return &Cars{l: l, v: v, cr: cr}
}

// ErrInvalidCarPath is an error message when the car path is not valid
var ErrInvalidCarPath = fmt.Errorf("Invalid Path, path should be /cars/[id]")

// GenericError is a generic error message returned by a server
type GenericError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"messages"`
}

// getCarId returns the car id from the URL
// Panic if cannot convert the id into an integer
// this should never happen as the router ensures that
// this is a valid number
func getCarId(r *http.Request) int {
	//parse the car id from the url
	vars := mux.Vars(r)

	//convert the id into an integer and return
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// should never happen
		panic(err)
	}
	return id
}
