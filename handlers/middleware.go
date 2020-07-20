package handlers

import (
	"context"
	"fmt"
	"github.com/CassioRoos/MicroseService/data"
	"net/http"
)

func (c Cars) MiddlewareValidateCar(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		car := &data.Car{}

		if err := data.FromJSON(car, r.Body); err != nil {
			c.l.Println("[ERROR] deserializing Car", err)
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		//Validate the car content before moving forward
		errs := c.v.Validate(car)
		if len(errs) != 0 {
			c.l.Println("[ERROR] validating Car", errs)
			http.Error(rw, fmt.Sprintf("Error reading the car: %s", errs), http.StatusBadRequest)
			return
		}
		// add the car to the context
		ctx := context.WithValue(r.Context(), KeyCar{}, *car)
		r = r.WithContext(ctx)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)

	})
}
