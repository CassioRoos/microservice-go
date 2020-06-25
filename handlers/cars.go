package handlers

import (
	"MicroseService/data"
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type (
	Cars struct {
		l *log.Logger
	}

	KeyCar struct{}
)

func NewCars(l *log.Logger) *Cars {
	return &Cars{l}
}

// func (c *Cars) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
// 	c.l.Println(r.Method)
// 	if r.Method == http.MethodGet {
// 		c.GetCars(rw, r)
// 		return
// 	}
// 	if r.Method == http.MethodPost {
// 		c.PostCar(rw, r)
// 		return
// 	}
// 	rw.WriteHeader(http.StatusMethodNotAllowed)
// }

func (c *Cars) GetCars(rw http.ResponseWriter, r *http.Request) {
	c.l.Println("Handle GET Cars")
	// Return the type CARS
	lc := data.GetCars()
	// FROM the type we can call the method to JSON, this way we can optimise our code
	if err := lc.ToJSON(rw); err != nil {
		http.Error(rw, "Unable to marshal JSON", http.StatusInternalServerError)
		return
	}
	// it's not necessary to call write, cuz we are already writing
	// rw.Write(data)
}

func (c *Cars) PostCar(rw http.ResponseWriter, r *http.Request) {
	c.l.Println("Handle POST ")

	car := r.Context().Value(KeyCar{}).(data.Car)
	data.AddCar(&car)
	c.l.Printf("Car %#v", car)
}

func (c *Cars) UpdateCar(rw http.ResponseWriter, r *http.Request) {
	// the ID is stored inside the mux
	vars := mux.Vars(r)
	// Get the param by name
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(rw, "Unable to convert the ID", http.StatusBadRequest)
		return
	}
	c.l.Println("Handle PUT Car", id)
	car := r.Context().Value(KeyCar{}).(data.Car)
	err = data.UpdateCar(id, &car)
	if err == data.ErrCarNotFound {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c Cars) MiddlewareValidateCar(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		car := data.Car{}

		if err := car.FromJSON(r.Body); err != nil {
			c.l.Println("Error deserializing Car", err)
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		// add the car to the context
		ctx := context.WithValue(r.Context(), KeyCar{}, car)
		r = r.WithContext(ctx)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)

	})
}
