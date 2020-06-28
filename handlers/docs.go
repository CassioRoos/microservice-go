// Package classificarion of cars
//
//  Documentation for Cars API
// Schemes: http
// BasePath: /
// version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
//
// - application/json
// swagger:meta
package handlers

import (
	"MicroseService/data"
)

// NOTE: types defined here are purely for documentation purposes
// these types are not used by anu of the handlers
// DO NOT USE THIS TYPES ANYWHERE BUT FOR DOCUMENTATION

// Generic error message returned as a string
// swagger:response errorResponse
type errorResponseWrapper struct {
	// Description of the error
	// in: body
	Body GenericError
}

// Validation errors defines as an array os strings
// swagger:response errorValidation
type erroValidationWrapper struct {
	// Collection of the errors
	// in: body
	Body ValidationError
}

// A list of cars
// swagger:response carsResponse
type carsResponseWrapper struct {
	// All current cars
	// in: body
	Body []data.Car
}

// Data structure representing a single car
// swagger:response carResponse
type carResponseWrapper struct {
	// A single car
	// in: body
	Body data.Car
}

// When there is no return
// swagger:response noContentResponse
type noContentResponseWrapper struct{}

// swagger:parameters updateCar createCar
type carParamsWrapper struct {
	// Car data structure to Update or Create.
	// Note: the Id field is ignored by both Create and Update operations
	// in: body
	// required: true
	Body data.Car
}

//swagger:parameters updateCar
type carIdParamsWrapper struct {
	// the Id of the car for which the operation relates
	// in: path
	// required: true
	Id int `json:"id"`
}