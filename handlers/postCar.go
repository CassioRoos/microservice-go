package handlers

import (
	"github.com/CassioRoos/MicroseService/data"
	"net/http"
)

// swagger:route POST /cars cars createCars
// Create a new car
//
// responses:
// 	201: carResponse
// 	422: errorValidation
// 	501: errorResponse

// Create handles POST requests to add new cars
func (c *Cars) PostCar(rw http.ResponseWriter, r *http.Request) {
	c.l.Println("Handle POST ")
	rw.WriteHeader(http.StatusCreated)
	rw.Header().Add("Content-Type", "application/json")
	car := r.Context().Value(KeyCar{}).(data.Car)
	data.AddCar(&car)
	c.l.Printf("Car %#v", car)
}
