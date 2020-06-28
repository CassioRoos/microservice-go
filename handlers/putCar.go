package handlers

import (
	"MicroseService/data"
	"net/http"
)


// swagger:route PUT /cars cars updateCar
// update the car details
//
// responses:
// 	200: noContentResponse
// 	404: errorResponse
// 	422: errorValidation

// Update handles PUT request to update car
func (c *Cars) UpdateCar(rw http.ResponseWriter, r *http.Request) {

	c.l.Println("Handle PUT Car")
	car := r.Context().Value(KeyCar{}).(data.Car)
	err := data.UpdateCar(car)
	if err == data.ErrCarNotFound {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
